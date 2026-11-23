package posttransform

import (
	"strings"
	"testing"

	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

const validRules = `
# Thresholds for CPU-like metrics.
transform_begin
transform_name       cpu thresholds
match_module         .*cpu_percent
match_extra_data     .*:metric:cpu
set_min_warning      80
set_max_critical     0
set_module_group     Rendimiento
set_tags             infra:cpu,team:ops
transform_end

transform_begin
transform_name       disk usage
match_module         .*disk_usage.*percent
set_min_warning      0.5
set_str_warning      high usage
transform_end
`

func TestParseValidBlocks(t *testing.T) {
	rules, err := ParseTransformBlocks(validRules)
	if err != nil {
		t.Fatalf("ParseTransformBlocks: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}

	if rules[0].Name != "cpu thresholds" {
		t.Errorf("rules[0].Name = %q", rules[0].Name)
	}
	if rules[0].Match["module"] != ".*cpu_percent" {
		t.Errorf("rules[0].Match[module] = %q", rules[0].Match["module"])
	}
	if rules[0].Match["extra_data"] != ".*:metric:cpu" {
		t.Errorf("rules[0].Match[extra_data] = %q", rules[0].Match["extra_data"])
	}
	if rules[0].Set["min_warning"] != 80 {
		t.Errorf("rules[0].Set[min_warning] = %v (%T), want int 80", rules[0].Set["min_warning"], rules[0].Set["min_warning"])
	}
	if rules[0].Set["max_critical"] != 0 {
		t.Errorf("rules[0].Set[max_critical] = %v, want int 0", rules[0].Set["max_critical"])
	}
	if rules[0].Set["module_group"] != "Rendimiento" {
		t.Errorf("rules[0].Set[module_group] = %v", rules[0].Set["module_group"])
	}
	if rules[0].Set["tags"] != "infra:cpu,team:ops" {
		t.Errorf("rules[0].Set[tags] = %v", rules[0].Set["tags"])
	}

	if v, ok := rules[1].Set["min_warning"].(float64); !ok || v != 0.5 {
		t.Errorf("rules[1].Set[min_warning] = %v (%T), want float64 0.5", rules[1].Set["min_warning"], rules[1].Set["min_warning"])
	}
	if rules[1].Set["str_warning"] != "high usage" {
		t.Errorf("rules[1].Set[str_warning] = %v", rules[1].Set["str_warning"])
	}
}

func TestParseErrors(t *testing.T) {
	cases := []struct {
		name     string
		text     string
		contains []string
	}{
		{
			"unknown set key",
			"transform_begin\nmatch_module a\nset_min_warnig 80\ntransform_end",
			[]string{"block 1", "min_warnig"},
		},
		{
			"unknown match key",
			"transform_begin\nmatch_modul a\nset_min_warning 1\ntransform_end",
			[]string{"unknown key", "match_modul"},
		},
		{
			"garbage line inside block",
			"transform_begin\nmatch_module a\nrandom junk\ntransform_end",
			[]string{"unknown key", "random"},
		},
		{
			"unparseable line inside block",
			"transform_begin\nmatch_module a\njunk!!\ntransform_end",
			[]string{"cannot parse line"},
		},
		{
			"junk outside block",
			"random junk outside\ntransform_begin\nmatch_module a\ntransform_end",
			[]string{"outside a transform block"},
		},
		{
			"unterminated block",
			"transform_begin\nmatch_module a\nset_min_warning 1\n",
			[]string{"missing transform_end"},
		},
		{
			"match-less rule",
			"transform_begin\nset_min_warning 1\ntransform_end",
			[]string{"no match_module or match_extra_data"},
		},
		{
			"invalid regex",
			"transform_begin\nmatch_module ([unclosed\nset_min_warning 1\ntransform_end",
			[]string{"block 1", "invalid regex"},
		},
		{
			"non-numeric threshold",
			"transform_begin\nmatch_module a\nset_min_warning abc\ntransform_end",
			[]string{"set_min_warning expects a number"},
		},
		{
			"int01 out of range",
			"transform_begin\nmatch_module a\nset_quiet 7\ntransform_end",
			[]string{"set_quiet expects 0 or 1"},
		},
		{
			"nested begin",
			"transform_begin\nmatch_module a\ntransform_begin\ntransform_end",
			[]string{"missing transform_end before a new transform_begin"},
		},
		{
			"end without begin",
			"transform_end\n",
			[]string{"without transform_begin"},
		},
	}

	for _, tc := range cases {
		_, err := ParseTransformBlocks(tc.text)
		if err == nil {
			t.Errorf("%s: expected error, got nil", tc.name)
			continue
		}
		for _, want := range tc.contains {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("%s: error %q does not contain %q", tc.name, err.Error(), want)
			}
		}
	}
}

func TestProtectedFields(t *testing.T) {
	for _, field := range []string{"name", "type", "data", "value", "extra_data", "timestamp", "status"} {
		text := "transform_begin\nmatch_module a\nset_" + field + " x\ntransform_end"
		_, err := ParseTransformBlocks(text)
		if err == nil {
			t.Errorf("protected field %q accepted", field)
			continue
		}
		if !strings.Contains(err.Error(), "protected") {
			t.Errorf("protected field %q: error %q does not say protected", field, err.Error())
		}
	}
}

func itemsFixture() []map[string]any {
	return []map[string]any{{
		"agent_data": map[string]any{"agent_name": "host_1"},
		"module_data": []map[string]any{
			{
				"name":       "host_1 cpu_percent",
				"type":       "generic_data",
				"data":       55,
				"extra_data": "nginx:metric:cpu",
			},
			{
				"name": "host_1 disk_usage_raw",
				"type": "generic_data",
				"data": 3,
			},
		},
	}}
}

func TestApplyMonitoringItems(t *testing.T) {
	rules, err := ParseTransformBlocks(validRules)
	if err != nil {
		t.Fatalf("ParseTransformBlocks: %v", err)
	}

	items := itemsFixture()
	changed := ApplyPostTransforms(items, rules)
	if changed != 1 {
		t.Fatalf("changed = %d, want 1", changed)
	}

	matched := items[0]["module_data"].([]map[string]any)[0]
	if matched["min_warning"] != 80 {
		t.Errorf("matched min_warning = %v", matched["min_warning"])
	}
	if matched["max_critical"] != 0 {
		t.Errorf("matched max_critical = %v", matched["max_critical"])
	}
	if matched["module_group"] != "Rendimiento" {
		t.Errorf("matched module_group = %v", matched["module_group"])
	}
	if matched["tags"] != "infra:cpu,team:ops" {
		t.Errorf("matched tags = %v", matched["tags"])
	}

	unmatched := items[0]["module_data"].([]map[string]any)[1]
	if _, ok := unmatched["min_warning"]; ok {
		t.Errorf("unmatched module was modified: %v", unmatched)
	}
}

func TestApplyLayeredLastMatchWinsPerField(t *testing.T) {
	layered := "transform_begin\nmatch_module .*cpu\nset_min_warning 10\nset_module_group Base\ntransform_end\n" +
		"transform_begin\nmatch_module .*cpu\nset_min_warning 20\ntransform_end"

	rules, err := ParseTransformBlocks(layered)
	if err != nil {
		t.Fatalf("ParseTransformBlocks: %v", err)
	}

	items := []map[string]any{{
		"agent_data":  map[string]any{"agent_name": "h2"},
		"module_data": []map[string]any{{"name": "h2 cpu", "type": "generic_data", "data": 1}},
	}}

	ApplyPostTransforms(items, rules)

	mod := items[0]["module_data"].([]map[string]any)[0]
	if mod["min_warning"] != 20 {
		t.Errorf("min_warning = %v, want 20 (last match wins)", mod["min_warning"])
	}
	if mod["module_group"] != "Base" {
		t.Errorf("module_group = %v, want Base (earlier block keeps its other fields)", mod["module_group"])
	}
}

func TestApplyMatchKeysAreAnded(t *testing.T) {
	rulesText := "transform_begin\nmatch_module .*cpu_percent\nmatch_extra_data .*:metric:cpu\nset_min_warning 99\ntransform_end"
	rules, err := ParseTransformBlocks(rulesText)
	if err != nil {
		t.Fatalf("ParseTransformBlocks: %v", err)
	}

	items := []map[string]any{{
		"agent_data": map[string]any{"agent_name": "h3"},
		"module_data": []map[string]any{{
			"name":       "h3 cpu_percent",
			"type":       "generic_data",
			"data":       1,
			"extra_data": "nginx:metric:memory",
		}},
	}}

	if changed := ApplyPostTransforms(items, rules); changed != 0 {
		t.Fatalf("changed = %d, want 0 (AND semantics)", changed)
	}
}

func TestApplyToModulesXML(t *testing.T) {
	rules, err := ParseTransformBlocks(validRules)
	if err != nil {
		t.Fatalf("ParseTransformBlocks: %v", err)
	}

	mod, err := pptmodule.New(pptmodule.Config{
		Name:      "x cpu_percent",
		Type:      "generic_data",
		Value:     "12",
		ExtraData: "nginx:metric:cpu",
	})
	if err != nil {
		t.Fatalf("module.New: %v", err)
	}

	modules := []pptmodule.Module{mod}
	if changed := ApplyToModules(modules, rules); changed != 1 {
		t.Fatalf("changed = %d, want 1", changed)
	}

	xml, err := pptmodule.XML(modules, nil)
	if err != nil {
		t.Fatalf("module.XML: %v", err)
	}
	rendered := string(xml)

	// The CPU rule matches by name and extra_data; its typed values must
	// reach the XML exactly as the Python engine would print them.
	for _, want := range []string{
		"<min_warning><![CDATA[80]]></min_warning>",
		"<max_critical><![CDATA[0]]></max_critical>",
		"<module_group>Rendimiento</module_group>",
		"<tags>infra:cpu,team:ops</tags>",
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("XML missing %q:\n%s", want, rendered)
		}
	}
}

func TestParseTolerantKeepsGoodRules(t *testing.T) {
	mixed := "transform_begin\nmatch_module .*good_one.*\nset_min_warning 42\ntransform_end\n\n" +
		"transform_begin\nmatch_module .*bad.*\nset_min_warnig 99\ntransform_end\n\n" +
		"transform_begin\nmatch_module .*good_two.*\nset_max_critical 7\ntransform_end"

	rules, errs := ParseTransformBlocksTolerant(mixed)
	if len(rules) != 2 {
		t.Fatalf("expected the two good rules, got %d (errs: %v)", len(rules), errs)
	}
	if rules[0].Set["min_warning"] != 42 {
		t.Errorf("rules[0].Set[min_warning] = %v", rules[0].Set["min_warning"])
	}
	if rules[1].Set["max_critical"] != 7 {
		t.Errorf("rules[1].Set[max_critical] = %v", rules[1].Set["max_critical"])
	}
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %v", errs)
	}
	if !strings.Contains(errs[0], "block 2") || !strings.Contains(errs[0], "min_warnig") {
		t.Errorf("error %q does not name block 2 / min_warnig", errs[0])
	}
}

func TestParseTolerantBadRegexSkipsOnlyItsRule(t *testing.T) {
	text := "transform_begin\nmatch_module .*good.*\nset_min_warning 42\ntransform_end\n\n" +
		"transform_begin\nmatch_module ([unclosed\nset_min_warning 1\ntransform_end"

	rules, errs := ParseTransformBlocksTolerant(text)
	if len(rules) != 1 {
		t.Fatalf("expected 1 good rule, got %d (errs: %v)", len(rules), errs)
	}
	if len(errs) != 1 || !strings.Contains(errs[0], "invalid regex") {
		t.Fatalf("expected one invalid-regex error, got %v", errs)
	}
}

func TestParseTolerantEmptyText(t *testing.T) {
	rules, errs := ParseTransformBlocksTolerant("")
	if len(rules) != 0 || len(errs) != 0 {
		t.Fatalf("expected empty results, got %v / %v", rules, errs)
	}
}

func TestApplyToModulesAlertTemplate(t *testing.T) {
	rulesText := "transform_begin\nmatch_module .*cpu_percent\nset_alert_template Critical mail\ntransform_end"
	rules, err := ParseTransformBlocks(rulesText)
	if err != nil {
		t.Fatalf("ParseTransformBlocks: %v", err)
	}

	mod, err := pptmodule.New(pptmodule.Config{
		Name:      "x cpu_percent",
		Type:      "generic_data",
		Value:     "12",
		ExtraData: "nginx:metric:cpu",
	})
	if err != nil {
		t.Fatalf("module.New: %v", err)
	}

	modules := []pptmodule.Module{mod}
	if changed := ApplyToModules(modules, rules); changed != 1 {
		t.Fatalf("changed = %d, want 1", changed)
	}
	if len(modules[0].Config.AlertTemplates) != 1 || modules[0].Config.AlertTemplates[0] != "Critical mail" {
		t.Fatalf("AlertTemplates = %v", modules[0].Config.AlertTemplates)
	}

	xml, err := pptmodule.XML(modules, nil)
	if err != nil {
		t.Fatalf("module.XML: %v", err)
	}
	if !strings.Contains(string(xml), "<alert_template><![CDATA[Critical mail]]></alert_template>") {
		t.Errorf("XML missing alert_template:\n%s", string(xml))
	}
}

func TestApplyToModulesLayered(t *testing.T) {
	layered := "transform_begin\nmatch_module .*cpu\nset_min_warning 10\nset_module_group Base\ntransform_end\n" +
		"transform_begin\nmatch_module .*cpu\nset_min_warning 20\nset_ff_timeout 30\ntransform_end"

	rules, err := ParseTransformBlocks(layered)
	if err != nil {
		t.Fatalf("ParseTransformBlocks: %v", err)
	}

	mod, err := pptmodule.New(pptmodule.Config{
		Name:      "node cpu_percent",
		Type:      "generic_data",
		Value:     "12",
		ExtraData: "nginx:metric:cpu",
	})
	if err != nil {
		t.Fatalf("module.New: %v", err)
	}

	modules := []pptmodule.Module{mod}
	ApplyToModules(modules, rules)

	cfg := modules[0].Config
	if cfg.MinWarning != "20" {
		t.Errorf("MinWarning = %q, want 20 (last match wins)", cfg.MinWarning)
	}
	if cfg.ModuleGroup != "Base" {
		t.Errorf("ModuleGroup = %q, want Base", cfg.ModuleGroup)
	}
	if cfg.FFTimeout != "30" {
		t.Errorf("FFTimeout = %q, want 30", cfg.FFTimeout)
	}
}
