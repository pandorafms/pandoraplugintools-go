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
)

type fixtureFile struct {
	Scenarios []scenario `json:"scenarios"`
}

type scenario struct {
	Name         string          `json:"name"`
	Agent        agentFixture    `json:"agent"`
	LogEncoding  string          `json:"log_encoding"`
	Modules      []moduleFixture `json:"modules"`
	LogModules   []logFixture    `json:"log_modules"`
	ImageModules []moduleFixture `json:"image_modules"`
	Discovery    discoFixture    `json:"discovery"`
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
	Name          string `json:"name"`
	AgentXML      string `json:"agent_xml"`
	DiscoveryJSON any    `json:"discovery_json"`
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

	for _, m := range s.Modules {
		mod, err := pptmodule.New(toModuleConfig(m))
		if err != nil {
			return "", fmt.Errorf("module.New: %w", err)
		}
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

	for _, entry := range s.Discovery.MonitoringDataAdds {
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

		results = append(results, result{Name: s.Name, AgentXML: agentXML, DiscoveryJSON: discoJSON})
	}

	out, err := json.Marshal(results)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(out))
}
