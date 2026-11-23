// Command go_harness builds agent XML and discovery JSON output for each
// scenario in parity/fixtures.json using pandoraplugintools-go, and prints
// one JSON object per scenario to stdout for compare.py to diff against the
// Python harness's output.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	pptagent "github.com/pandorafms/pandoraplugintools-go/pkg/agent"
	pptdiscovery "github.com/pandorafms/pandoraplugintools-go/pkg/discovery"
	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
	pptposttransform "github.com/pandorafms/pandoraplugintools-go/pkg/posttransform"
)

type fixtureFile struct {
	Scenarios []scenario `json:"scenarios"`
}

type scenario struct {
	Name           string                `json:"name"`
	Agent          agentFixture          `json:"agent"`
	LogEncoding    string                `json:"log_encoding"`
	Modules        []moduleFixture       `json:"modules"`
	LogModules     []logFixture          `json:"log_modules"`
	ImageModules   []moduleFixture       `json:"image_modules"`
	Discovery      discoFixture          `json:"discovery"`
	PostTransforms *postTransformFixture `json:"post_transforms"`
}

type postTransformFixture struct {
	Rules string `json:"rules"`
}

type agentFixture struct {
	AgentName       string `json:"agent_name"`
	AgentAlias      string `json:"agent_alias"`
	ParentAgentName string `json:"parent_agent_name"`
	Description     string `json:"description"`
	Version         string `json:"version"`
	OSName          string `json:"os_name"`
	OSVersion       string `json:"os_version"`
	Timestamp       string `json:"timestamp"`
	Address         string `json:"address"`
	Group           string `json:"group"`
	Interval        int    `json:"interval"`
	AgentMode       string `json:"agent_mode"`
}

type dataPointFixture struct {
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}

type moduleFixture struct {
	Name           string             `json:"name"`
	Type           string             `json:"type"`
	Value          string             `json:"value"`
	Unit           string             `json:"unit"`
	Tags           string             `json:"tags"`
	ModuleGroup    string             `json:"module_group"`
	ModuleParent   string             `json:"module_parent"`
	CustomID       string             `json:"custom_id"`
	ExtraData      string             `json:"extra_data"`
	MinWarning     string             `json:"min_warning"`
	MaxWarning     string             `json:"max_warning"`
	MinCritical    string             `json:"min_critical"`
	MaxCritical    string             `json:"max_critical"`
	Min            string             `json:"min"`
	Max            string             `json:"max"`
	DataList       []dataPointFixture `json:"data_list"`
	AlertTemplates []string           `json:"alert_templates"`
}

type logFixture struct {
	Source string `json:"source"`
	Value  string `json:"value"`
}

type kvFixture struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type discoFixture struct {
	ErrorLevel         int              `json:"error_level"`
	SummaryValues      []kvFixture      `json:"summary_values"`
	SummaryAdds        []kvFixture      `json:"summary_adds"`
	InfoSet            *string          `json:"info_set"`
	InfoAdds           []string         `json:"info_adds"`
	MonitoringDataAdds []map[string]any `json:"monitoring_data_adds"`
}

type result struct {
	Name                  string `json:"name"`
	AgentXML              string `json:"agent_xml"`
	DiscoveryJSON         any    `json:"discovery_json"`
	ParsedRules           any    `json:"parsed_rules,omitempty"`
	MonitoringItemsAfter  any    `json:"monitoring_items_after,omitempty"`
	PostTransformErrorCnt any    `json:"post_transform_error_count,omitempty"`
}

func toModuleConfig(m moduleFixture) pptmodule.Config {
	cfg := pptmodule.Config{
		Name:           m.Name,
		Type:           m.Type,
		Value:          m.Value,
		Unit:           m.Unit,
		Tags:           m.Tags,
		ModuleGroup:    m.ModuleGroup,
		ModuleParent:   m.ModuleParent,
		CustomID:       m.CustomID,
		ExtraData:      m.ExtraData,
		MinWarning:     m.MinWarning,
		MaxWarning:     m.MaxWarning,
		MinCritical:    m.MinCritical,
		MaxCritical:    m.MaxCritical,
		Min:            m.Min,
		Max:            m.Max,
		AlertTemplates: m.AlertTemplates,
	}

	for _, dp := range m.DataList {
		cfg.DataList = append(cfg.DataList, pptmodule.DataPoint{Value: dp.Value, Timestamp: dp.Timestamp})
	}

	return cfg
}

func parseScenarioRules(s scenario) ([]pptposttransform.Rule, []string) {
	if s.PostTransforms == nil {
		return nil, nil
	}
	// Tolerant parse: one bad rule is skipped, the good rules still apply.
	return pptposttransform.ParseTransformBlocksTolerant(s.PostTransforms.Rules)
}

func buildAgentXML(s scenario) (string, error) {
	a, err := pptagent.New(pptagent.Config{
		AgentName:       s.Agent.AgentName,
		AgentAlias:      s.Agent.AgentAlias,
		ParentAgentName: s.Agent.ParentAgentName,
		Description:     s.Agent.Description,
		Version:         s.Agent.Version,
		OSName:          s.Agent.OSName,
		OSVersion:       s.Agent.OSVersion,
		Timestamp:       s.Agent.Timestamp,
		Address:         s.Agent.Address,
		Group:           s.Agent.Group,
		Interval:        s.Agent.Interval,
		AgentMode:       s.Agent.AgentMode,
	})
	if err != nil {
		return "", fmt.Errorf("agent.New: %w", err)
	}

	modules := make([]pptmodule.Module, 0, len(s.Modules))
	for _, m := range s.Modules {
		mod, err := pptmodule.New(toModuleConfig(m))
		if err != nil {
			return "", fmt.Errorf("module.New: %w", err)
		}
		modules = append(modules, mod)
	}

	// Post-transforms run over the complete module list right before the
	// final serialization (ticket #22040 XML path).
	rules, _ := parseScenarioRules(s)
	if rules != nil {
		pptposttransform.ApplyToModules(modules, rules)
	}

	for _, mod := range modules {
		if err := a.AddModule(mod); err != nil {
			return "", fmt.Errorf("AddModule: %w", err)
		}
	}

	for _, lm := range s.LogModules {
		mod, err := pptmodule.NewLog(pptmodule.LogConfig{Source: lm.Source, Value: lm.Value})
		if err != nil {
			return "", fmt.Errorf("module.NewLog: %w", err)
		}
		if err := a.AddLogModule(mod); err != nil {
			return "", fmt.Errorf("AddLogModule: %w", err)
		}
	}

	for _, m := range s.ImageModules {
		mod, err := pptmodule.New(toModuleConfig(m))
		if err != nil {
			return "", fmt.Errorf("image module.New: %w", err)
		}
		if err := a.AddImageModule(mod); err != nil {
			return "", fmt.Errorf("AddImageModule: %w", err)
		}
	}

	body, err := a.XMLWithOptions(pptagent.XMLOptions{LogEncoding: s.LogEncoding})
	if err != nil {
		return "", fmt.Errorf("Agent.XMLWithOptions: %w", err)
	}

	return string(body), nil
}

func buildDiscoveryJSON(s scenario) (any, error) {
	d := pptdiscovery.New()
	d.SetErrorLevel(s.Discovery.ErrorLevel)

	for _, kv := range s.Discovery.SummaryValues {
		d.SetSummaryValue(kv.Key, normalizeJSONNumber(kv.Value))
	}

	for _, kv := range s.Discovery.SummaryAdds {
		if err := d.AddSummaryValue(kv.Key, normalizeJSONNumber(kv.Value)); err != nil {
			return nil, fmt.Errorf("AddSummaryValue: %w", err)
		}
	}

	if s.Discovery.InfoSet != nil {
		d.SetInfo(*s.Discovery.InfoSet)
	}

	for _, text := range s.Discovery.InfoAdds {
		d.AddInfo(text)
	}

	// Post-transforms run over the accumulated monitoring items before
	// OutputJSON serializes them (ticket #22040 discovery JSON path).
	entries := s.Discovery.MonitoringDataAdds
	rules, _ := parseScenarioRules(s)
	if rules != nil {
		pptposttransform.ApplyPostTransforms(entries, rules)
	}

	for _, entry := range entries {
		d.AddMonitoringData(entry)
	}

	out, err := d.OutputJSON()
	if err != nil {
		return nil, fmt.Errorf("OutputJSON: %w", err)
	}

	var parsed any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		return nil, fmt.Errorf("unmarshal discovery output: %w", err)
	}

	return parsed, nil
}

// normalizeJSONNumber converts encoding/json's default float64 decoding into
// an int when the fixture value has no fractional part, so AddSummaryValue's
// int/float64/string type switch matches what the fixture author intended
// (JSON has no separate int type, everything decodes as float64 otherwise).
func normalizeJSONNumber(v any) any {
	f, ok := v.(float64)
	if !ok {
		return v
	}
	if f == float64(int(f)) {
		return int(f)
	}
	return f
}

// rawScenario mirrors the fixture JSON, keeping the raw agent/module maps
// the Python harness also reads, so the monitoring items built for the
// "/monitoring" path are identical on both sides.
type rawScenario struct {
	Name    string           `json:"name"`
	Agent   map[string]any   `json:"agent"`
	Modules []map[string]any `json:"modules"`
}

// translateRawModule applies the same fixture-to-library translation as the
// Python harness's to_python_module_dict: alert_templates -> alert and
// data_list -> the value list print_module detects as datalist mode.
func translateRawModule(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}

	if v, ok := out["alert_templates"]; ok {
		delete(out, "alert_templates")
		out["alert"] = v
	}

	if v, ok := out["data_list"]; ok {
		delete(out, "data_list")
		points, _ := v.([]any)
		list := make([]any, 0, len(points))
		for _, p := range points {
			pm, _ := p.(map[string]any)
			entry := map[string]any{"value": pm["value"]}
			if ts, ok := pm["timestamp"]; ok {
				if s, isStr := ts.(string); !isStr || s != "" {
					entry["timestamp"] = ts
				}
			}
			list = append(list, entry)
		}
		out["value"] = list
	}

	return out
}

// buildPostTransformOutputs returns the parsed rules (name/match/set, typed
// values), the transformed monitoring items and the count of rule errors the
// tolerant parser skipped (the error TEXT is per-language, so only the count
// is compared), for scenarios that carry post_transforms; nil, nil, nil,
// nil otherwise. This proves parser parity and the "/monitoring" JSON path
// on top of the agent XML and discovery JSON paths already exercised by the
// scenario itself.
func buildPostTransformOutputs(s scenario, raw rawScenario) (any, any, any, error) {
	if s.PostTransforms == nil {
		return nil, nil, nil, nil
	}

	rules, ptErrs := parseScenarioRules(s)

	modules := make([]map[string]any, 0, len(raw.Modules))
	for _, m := range raw.Modules {
		modules = append(modules, translateRawModule(m))
	}
	items := []map[string]any{{
		"agent_data":  raw.Agent,
		"module_data": modules,
	}}
	pptposttransform.ApplyPostTransforms(items, rules)

	var parsedRules any
	if b, err := json.Marshal(rules); err == nil {
		if err := json.Unmarshal(b, &parsedRules); err != nil {
			return nil, nil, nil, fmt.Errorf("unmarshal parsed rules: %w", err)
		}
	}

	var itemsAfter any
	if b, err := json.Marshal(items); err == nil {
		if err := json.Unmarshal(b, &itemsAfter); err != nil {
			return nil, nil, nil, fmt.Errorf("unmarshal monitoring items: %w", err)
		}
	}

	return parsedRules, itemsAfter, len(ptErrs), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go_harness <fixtures.json>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var fixtures fixtureFile
	if err := json.Unmarshal(data, &fixtures); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Shadow unmarshal keeping the raw agent/module maps, used by the
	// post-transform "/monitoring" path so both harnesses build the items
	// from the same JSON the fixture author wrote.
	var rawFile struct {
		Scenarios []rawScenario `json:"scenarios"`
	}
	if err := json.Unmarshal(data, &rawFile); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rawByName := make(map[string]rawScenario, len(rawFile.Scenarios))
	for _, rs := range rawFile.Scenarios {
		rawByName[rs.Name] = rs
	}

	results := make([]result, 0, len(fixtures.Scenarios))
	for _, s := range fixtures.Scenarios {
		agentXML, err := buildAgentXML(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "scenario %q: %v\n", s.Name, err)
			os.Exit(1)
		}

		discoJSON, err := buildDiscoveryJSON(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "scenario %q: %v\n", s.Name, err)
			os.Exit(1)
		}

		parsedRules, itemsAfter, ptErrCount, err := buildPostTransformOutputs(s, rawByName[s.Name])
		if err != nil {
			fmt.Fprintf(os.Stderr, "scenario %q: %v\n", s.Name, err)
			os.Exit(1)
		}

		results = append(results, result{
			Name:                  s.Name,
			AgentXML:              agentXML,
			DiscoveryJSON:         discoJSON,
			ParsedRules:           parsedRules,
			MonitoringItemsAfter:  itemsAfter,
			PostTransformErrorCnt: ptErrCount,
		})
	}

	out, err := json.Marshal(results)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(out))
}
