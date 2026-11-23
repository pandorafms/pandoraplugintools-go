// Package posttransform ports pandoraPlugintoolsBasic's post_transforms.py
// (ticket #22040): ordered regex rules, described as plain-text blocks, that
// merge whitelisted module configuration fields into the collected structure
// right before it is serialized, so the same rules work for every output path
// of both libraries:
//
//   - XML (tentacle): ApplyToModules over the []module.Module list before
//     agent rendering.
//   - Discovery JSON (Discovery.OutputJSON): ApplyPostTransforms over the
//     []map[string]any monitoring items before SetMonitoringData.
//   - "/monitoring" JSON (monitoring.PayloadJSON): the same
//     ApplyPostTransforms over the items before building the payload.
//
// Block grammar (one rule per block):
//
//	transform_begin
//	transform_name       <optional label>
//	match_module         <regex against the module name>
//	match_extra_data     <regex against the module extra_data>
//	set_<field>          <typed value from the whitelist>
//	transform_end
//
// Semantics are shared with the Python engine: comments (#) and blank lines
// are ignored; match keys inside a block are AND-ed and a rule needs at least
// one; rules apply in order and for each field the last matching rule wins;
// 'name', 'type', 'data', 'value', 'extra_data', 'timestamp' and 'status' are
// protected; threshold ranges are half-open [min, max) with max_warning=0 /
// max_critical=0 meaning unbounded. Unknown keys, invalid regexes and
// mistyped values are reported as rule errors naming the block and the key:
// the tolerant API skips the offending rule and keeps the rest (one bad
// rule never throws away the monitoring run), while the strict API returns
// them as an error instead.
//
// Matching uses Go regexp (RE2) with MatchString's unanchored semantics,
// which mirrors Python's re.search. Rules must therefore stay
// RE2-compatible (no backreferences or lookaround) to behave identically in
// both engines.
package posttransform

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

// Rule is one parsed transform block. The JSON tags mirror the Python
// engine's rule dicts ({"name", "match", "set"}) so parity fixtures compare
// directly. Set values are typed: int for integral numbers and int01 fields,
// float64 for fractional numbers, string otherwise.
type Rule struct {
	Name  string            `json:"name"`
	Match map[string]string `json:"match"`
	Set   map[string]any    `json:"set"`
}

// fieldSpec describes one whitelisted set_<field> key.
type fieldSpec struct {
	// Kind is one of "number", "int", "int01", "str", "alert".
	Kind string
	// ConfigField is the module.Config field the value maps to on the XML
	// path ("" for "alert", handled by appending to AlertTemplates).
	ConfigField string
}

// whitelist mirrors the Python engine's WHITELIST: exactly the module fields
// both libraries already template and serialize, i.e. the vocabulary the
// Discovery ingestion accepts.
var whitelist = map[string]fieldSpec{
	"min_warning":           {"number", "MinWarning"},
	"max_warning":           {"number", "MaxWarning"},
	"min_critical":          {"number", "MinCritical"},
	"max_critical":          {"number", "MaxCritical"},
	"min_warning_forced":    {"number", "MinWarningForced"},
	"max_warning_forced":    {"number", "MaxWarningForced"},
	"min_critical_forced":   {"number", "MinCriticalForced"},
	"max_critical_forced":   {"number", "MaxCriticalForced"},
	"min":                   {"number", "Min"},
	"max":                   {"number", "Max"},
	"post_process":          {"number", "PostProcess"},
	"str_warning":           {"str", "StrWarning"},
	"str_critical":          {"str", "StrCritical"},
	"str_warning_forced":    {"str", "StrWarningForced"},
	"str_critical_forced":   {"str", "StrCriticalForced"},
	"critical_inverse":      {"int01", "CriticalInverse"},
	"warning_inverse":       {"int01", "WarningInverse"},
	"quiet":                 {"int01", "Quiet"},
	"disabled":              {"int01", "Disabled"},
	"each_ff":               {"int01", "EachFF"},
	"min_ff_event":          {"int", "MinFFEvent"},
	"min_ff_event_normal":   {"int", "MinFFEventNormal"},
	"min_ff_event_warning":  {"int", "MinFFEventWarning"},
	"min_ff_event_critical": {"int", "MinFFEventCritical"},
	"ff_type":               {"int", "FFType"},
	"ff_timeout":            {"int", "FFTimeout"},
	"module_ff_interval":    {"int", "ModuleFFInterval"},
	"interval":              {"int", "Interval"},
	"crontab":               {"str", "CronTab"},
	"tags":                  {"str", "Tags"},
	"module_group":          {"str", "ModuleGroup"},
	"module_parent":         {"str", "ModuleParent"},
	"module_parent_unlink":  {"str", "ModuleParentUnlink"},
	"custom_id":             {"str", "CustomID"},
	"description":           {"str", "Description"},
	"critical_instructions": {"str", "CriticalInstructions"},
	"warning_instructions":  {"str", "WarningInstructions"},
	"unknown_instructions":  {"str", "UnknownInstructions"},
	"alert_template":        {"alert", "AlertTemplates"},
}

// protectedFields can never be written by a rule: they define module
// identity and the collected value.
var protectedFields = map[string]bool{
	"name":       true,
	"type":       true,
	"data":       true,
	"value":      true,
	"extra_data": true,
	"timestamp":  true,
	"status":     true,
}

var lineRE = regexp.MustCompile(`^([a-z_][a-z0-9_]*)\s+(.+)$`)

const (
	matchKeyModule    = "match_module"
	matchKeyExtraData = "match_extra_data"
)

func cleanLine(s string) string {
	s = strings.ReplaceAll(s, "\u200b", "")
	s = strings.ReplaceAll(s, "\ufeff", "")
	return strings.TrimSpace(s)
}

// parseBlocksTolerant is the tolerant parse core: it parses every block,
// skipping the invalid RULES while keeping the valid ones, and collects one
// error string per problem, each naming its block and line. Rule granularity:
// one bad key invalidates its whole rule, because applying half of an
// intended configuration is worse than skipping it.
func parseBlocksTolerant(rulesText string) ([]Rule, []string) {
	var rules []Rule
	var errs []string
	var current *Rule
	var currentErrs []string
	blockNo := 0

	for i, rawLine := range strings.Split(rulesText, "\n") {
		lineno := i + 1
		line := cleanLine(rawLine)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if line == "transform_begin" {
			if current != nil {
				errs = append(errs, fmt.Sprintf(
					"transform block %d, line %d: missing transform_end before a new transform_begin; rule skipped",
					blockNo, lineno))
			}
			blockNo++
			current = &Rule{Match: map[string]string{}, Set: map[string]any{}}
			currentErrs = nil
			continue
		}

		if line == "transform_end" {
			if current == nil {
				errs = append(errs, fmt.Sprintf(
					"line %d: transform_end without transform_begin", lineno))
				continue
			}

			if len(current.Match) == 0 {
				currentErrs = append(currentErrs, fmt.Sprintf(
					"transform block %d: rule has no match_module or match_extra_data key; rule skipped",
					blockNo))
			}

			for key, pattern := range current.Match {
				if _, err := regexp.Compile(pattern); err != nil {
					currentErrs = append(currentErrs, fmt.Sprintf(
						"transform block %d: invalid regex in match key %q: %v; rule skipped",
						blockNo, key, err))
				}
			}

			if len(currentErrs) > 0 {
				errs = append(errs, currentErrs...)
			} else {
				rules = append(rules, *current)
			}

			current = nil
			currentErrs = nil
			continue
		}

		if current == nil {
			errs = append(errs, fmt.Sprintf(
				"line %d: content outside a transform block: %q", lineno, line))
			continue
		}

		parsed := lineRE.FindStringSubmatch(line)
		if parsed == nil {
			currentErrs = append(currentErrs, fmt.Sprintf(
				"transform block %d, line %d: cannot parse line: %q; rule skipped",
				blockNo, lineno, line))
			continue
		}

		key := parsed[1]
		value := strings.TrimSpace(parsed[2])

		switch {
		case key == "transform_name":
			current.Name = value
		case key == matchKeyModule:
			current.Match["module"] = value
		case key == matchKeyExtraData:
			current.Match["extra_data"] = value
		case strings.HasPrefix(key, "set_"):
			field := strings.TrimPrefix(key, "set_")
			if protectedFields[field] {
				currentErrs = append(currentErrs, fmt.Sprintf(
					"transform block %d, line %d: field %q is protected and cannot be set; rule skipped",
					blockNo, lineno, field))
				continue
			}
			spec, ok := whitelist[field]
			if !ok {
				currentErrs = append(currentErrs, fmt.Sprintf(
					"transform block %d, line %d: unknown set field %q (not in the whitelist); rule skipped",
					blockNo, lineno, field))
				continue
			}
			coerced, err := coerceValue(field, spec.Kind, value, blockNo, lineno)
			if err != nil {
				currentErrs = append(currentErrs, fmt.Sprintf("%v; rule skipped", err))
				continue
			}
			current.Set[field] = coerced
		default:
			currentErrs = append(currentErrs, fmt.Sprintf(
				"transform block %d, line %d: unknown key %q; rule skipped",
				blockNo, lineno, key))
		}
	}

	if current != nil {
		errs = append(errs, fmt.Sprintf(
			"transform block %d: missing transform_end at end of input; rule skipped",
			blockNo))
	}

	return rules, errs
}

// ParseTransformBlocks parses rule blocks into an ordered []Rule.
//
// STRICT variant: it returns an error naming the invalid rules instead of
// returning partial results. Use ParseTransformBlocksTolerant when one bad
// rule must never stop the monitoring run.
func ParseTransformBlocks(rulesText string) ([]Rule, error) {
	rules, errs := parseBlocksTolerant(rulesText)
	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}
	return rules, nil
}

// ParseTransformBlocksTolerant is the tolerant variant for plugins: one bad
// rule never throws away the whole monitoring run. Valid rules apply
// normally and each error names its block and line so the task can report
// it (through the discovery info plus an error level, and stderr for
// manual runs) while continuing.
func ParseTransformBlocksTolerant(rulesText string) ([]Rule, []string) {
	return parseBlocksTolerant(rulesText)
}

func coerceValue(field, kind, raw string, blockNo, lineno int) (any, error) {
	where := fmt.Sprintf("transform block %d, line %d", blockNo, lineno)

	switch kind {
	case "number":
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return nil, fmt.Errorf("%s: set_%s expects a number, got %q", where, field, raw)
		}
		if !math.IsInf(f, 0) && f == math.Trunc(f) {
			return int(f), nil
		}
		return f, nil
	case "int", "int01":
		n, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Errorf("%s: set_%s expects an integer, got %q", where, field, raw)
		}
		if kind == "int01" && n != 0 && n != 1 {
			return nil, fmt.Errorf("%s: set_%s expects 0 or 1, got %q", where, field, raw)
		}
		return n, nil
	default: // "str", "alert"
		return raw, nil
	}
}

// configValueString renders a typed Set value the way Python's str() would,
// so the XML path prints the same text both engines produce.
func configValueString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	default:
		return fmt.Sprintf("%v", value)
	}
}

// matchString reports whether pattern matches s with the same semantics as
// Python's re.search (unanchored); an invalid pattern simply never matches at
// apply time (parse already rejected invalid regexes).
func matchString(pattern, s string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(s)
}

func mapString(v any) string {
	switch s := v.(type) {
	case nil:
		return ""
	case string:
		return s
	default:
		return fmt.Sprintf("%v", s)
	}
}

func ruleMatches(moduleName, extraData string, rule Rule) bool {
	if pattern, ok := rule.Match["module"]; ok {
		if !matchString(pattern, moduleName) {
			return false
		}
	}
	if pattern, ok := rule.Match["extra_data"]; ok {
		if !matchString(pattern, extraData) {
			return false
		}
	}
	return true
}

func applyRuleToMap(module map[string]any, rule Rule) bool {
	touched := false
	for field, value := range rule.Set {
		if whitelist[field].Kind == "alert" {
			alerts, _ := module["alert"].([]any)
			module["alert"] = append(alerts, value)
		} else {
			module[field] = value
		}
		touched = true
	}
	return touched
}

func entryModules(entry map[string]any) []map[string]any {
	switch v := entry["module_data"].(type) {
	case []map[string]any:
		return v
	case []any:
		mods := make([]map[string]any, 0, len(v))
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				mods = append(mods, m)
			}
		}
		return mods
	default:
		return []map[string]any{entry}
	}
}

// ApplyPostTransforms applies parsed rules over collected monitoring data,
// in place, before serializing it to its final format. data accepts the
// monitoring items list ([{"agent_data": ..., "module_data": [...]}], used by
// the discovery JSON and "/monitoring" JSON paths) or a plain list of module
// maps. It returns the number of modules that received at least one field.
func ApplyPostTransforms(data []map[string]any, rules []Rule) int {
	changed := 0

	for _, entry := range data {
		for _, module := range entryModules(entry) {
			if module == nil {
				continue
			}
			touched := false
			for _, rule := range rules {
				if ruleMatches(mapString(module["name"]), mapString(module["extra_data"]), rule) &&
					applyRuleToMap(module, rule) {
					touched = true
				}
			}
			if touched {
				changed++
			}
		}
	}

	return changed
}

func setConfigField(cfg *pptmodule.Config, field, value string) {
	v := reflect.ValueOf(cfg).Elem().FieldByName(field)
	if v.IsValid() && v.Kind() == reflect.String {
		v.SetString(value)
	}
}

func applyRuleToConfig(cfg *pptmodule.Config, rule Rule) bool {
	touched := false
	for field, value := range rule.Set {
		spec := whitelist[field]
		if spec.Kind == "alert" {
			cfg.AlertTemplates = append(cfg.AlertTemplates, configValueString(value))
		} else {
			setConfigField(cfg, spec.ConfigField, configValueString(value))
		}
		touched = true
	}
	return touched
}

// ApplyToModules applies parsed rules over a []module.Module list, in place,
// before the XML/tentacle path renders it (agent rendering or module.XML).
// It returns the number of modules that received at least one field.
func ApplyToModules(modules []pptmodule.Module, rules []Rule) int {
	changed := 0

	for i := range modules {
		cfg := &modules[i].Config
		touched := false
		for _, rule := range rules {
			if ruleMatches(cfg.Name, cfg.ExtraData, rule) && applyRuleToConfig(cfg, rule) {
				touched = true
			}
		}
		if touched {
			changed++
		}
	}

	return changed
}
