package agent_test

import (
	"strings"
	"testing"

	pptagent "github.com/pandorafms/pandoraplugintools-go/pkg/agent"
	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

func TestAgentXMLIncludesScalarModule(t *testing.T) {
	ag, err := pptagent.New(pptagent.Config{AgentName: "agent-123", AgentAlias: "WIN-SERV"})
	if err != nil {
		t.Fatalf("expected agent to be created, got error: %v", err)
	}

	mod, err := pptmodule.New(pptmodule.Config{
		Name:                 "CPU usage",
		Type:                 "generic_data",
		Data:                 "10",
		Desc:                 "Percentage of CPU utilization",
		Unit:                 "%",
		Interval:             "300",
		Tags:                 "prod,linux",
		ModuleGroup:          "system",
		ModuleParent:         "host-health",
		MinWarning:           "60",
		MinWarningForced:     "65",
		MaxWarning:           "80",
		MaxWarningForced:     "85",
		MinCritical:          "20",
		MinCriticalForced:    "15",
		MaxCritical:          "95",
		MaxCriticalForced:    "98",
		StrWarning:           "warn",
		StrWarningForced:     "warn-forced",
		StrCritical:          "crit",
		StrCriticalForced:    "crit-forced",
		CriticalInverse:      "1",
		WarningInverse:       "0",
		Min:                  "0",
		Max:                  "100",
		PostProcess:          "scale(0.1)",
		Disabled:             "0",
		MinFFEvent:           "2",
		Status:               "normal",
		Timestamp:            "2026-05-22T08:00:00Z",
		CustomID:             "cpu-usage",
		CriticalInstructions: "Check sustained CPU pressure",
		WarningInstructions:  "Watch sustained load",
		UnknownInstructions:  "Check plugin execution",
		Quiet:                "1",
		ModuleFFInterval:     "60",
		CronTab:              "*/5 * * * *",
		MinFFEventNormal:     "1",
		MinFFEventWarning:    "2",
		MinFFEventCritical:   "3",
		FFType:               "jump",
		FFTimeout:            "30",
		EachFF:               "1",
		ModuleParentUnlink:   "0",
		ExtraData:            "sample=1",
		Alert:                []string{"cpu-warning", "cpu-critical"},
	})
	if err != nil {
		t.Fatalf("expected module to be created, got error: %v", err)
	}

	if err := ag.AddModule(mod); err != nil {
		t.Fatalf("expected module to be added, got error: %v", err)
	}

	xmlData, err := ag.XML()
	if err != nil {
		t.Fatalf("expected xml to be generated, got error: %v", err)
	}

	xmlString := string(xmlData)
	checks := []string{
		"<agent_data",
		`agent_name="agent-123"`,
		"<module>",
		"<![CDATA[CPU usage]]>",
		"<![CDATA[10]]>",
		"<description><![CDATA[Percentage of CPU utilization]]></description>",
		"<module_interval><![CDATA[300]]></module_interval>",
		"<tags>prod,linux</tags>",
		"<module_group>system</module_group>",
		"<module_parent>host-health</module_parent>",
		"<min_warning><![CDATA[60]]></min_warning>",
		"<min_warning_forced><![CDATA[65]]></min_warning_forced>",
		"<max_warning><![CDATA[80]]></max_warning>",
		"<max_warning_forced><![CDATA[85]]></max_warning_forced>",
		"<min_critical><![CDATA[20]]></min_critical>",
		"<min_critical_forced><![CDATA[15]]></min_critical_forced>",
		"<max_critical><![CDATA[95]]></max_critical>",
		"<max_critical_forced><![CDATA[98]]></max_critical_forced>",
		"<str_warning><![CDATA[warn]]></str_warning>",
		"<str_warning_forced><![CDATA[warn-forced]]></str_warning_forced>",
		"<str_critical><![CDATA[crit]]></str_critical>",
		"<str_critical_forced><![CDATA[crit-forced]]></str_critical_forced>",
		"<critical_inverse><![CDATA[1]]></critical_inverse>",
		"<warning_inverse><![CDATA[0]]></warning_inverse>",
		"<min><![CDATA[0]]></min>",
		"<max><![CDATA[100]]></max>",
		"<post_process><![CDATA[scale(0.1)]]></post_process>",
		"<disabled><![CDATA[0]]></disabled>",
		"<min_ff_event><![CDATA[2]]></min_ff_event>",
		"<status><![CDATA[normal]]></status>",
		"<timestamp><![CDATA[2026-05-22T08:00:00Z]]></timestamp>",
		"<custom_id><![CDATA[cpu-usage]]></custom_id>",
		"<critical_instructions><![CDATA[Check sustained CPU pressure]]></critical_instructions>",
		"<warning_instructions><![CDATA[Watch sustained load]]></warning_instructions>",
		"<unknown_instructions><![CDATA[Check plugin execution]]></unknown_instructions>",
		"<quiet><![CDATA[1]]></quiet>",
		"<module_ff_interval><![CDATA[60]]></module_ff_interval>",
		"<crontab><![CDATA[*/5 * * * *]]></crontab>",
		"<min_ff_event_normal><![CDATA[1]]></min_ff_event_normal>",
		"<min_ff_event_warning><![CDATA[2]]></min_ff_event_warning>",
		"<min_ff_event_critical><![CDATA[3]]></min_ff_event_critical>",
		"<ff_type><![CDATA[jump]]></ff_type>",
		"<ff_timeout><![CDATA[30]]></ff_timeout>",
		"<each_ff><![CDATA[1]]></each_ff>",
		"<module_parent_unlink><![CDATA[0]]></module_parent_unlink>",
		"<extra_data><![CDATA[sample=1]]></extra_data>",
		"<alert_template><![CDATA[cpu-warning]]></alert_template>",
		"<alert_template><![CDATA[cpu-critical]]></alert_template>",
	}

	for _, check := range checks {
		if !strings.Contains(xmlString, check) {
			t.Fatalf("expected XML to contain %q, got:\n%s", check, xmlString)
		}
	}
}

func TestAgentXMLWithOptionsIncludesDataListAndLogModules(t *testing.T) {
	ag, err := pptagent.New(pptagent.Config{AgentName: "agent-logs", AgentAlias: "agent-logs"})
	if err != nil {
		t.Fatalf("expected agent to be created, got error: %v", err)
	}

	series, err := pptmodule.New(pptmodule.Config{
		Name: "Process count",
		Type: "generic_data",
		DataList: []pptmodule.DataPoint{
			{Value: "10", Timestamp: "2026-05-22T10:00:00Z"},
			{Value: "12"},
		},
	})
	if err != nil {
		t.Fatalf("expected datalist module to be created, got error: %v", err)
	}

	logModule, err := pptmodule.NewLog(pptmodule.LogConfig{
		Source: "application.log",
		Value:  "Service restarted",
	})
	if err != nil {
		t.Fatalf("expected log module to be created, got error: %v", err)
	}

	if err := ag.AddModule(series); err != nil {
		t.Fatalf("expected datalist module to be added, got error: %v", err)
	}

	if err := ag.AddLogModule(logModule); err != nil {
		t.Fatalf("expected log module to be added, got error: %v", err)
	}

	xmlData, err := ag.XMLWithOptions(pptagent.XMLOptions{LogEncoding: "utf-8"})
	if err != nil {
		t.Fatalf("expected xml to be generated, got error: %v", err)
	}

	xmlString := string(xmlData)
	checks := []string{
		"<datalist>",
		"<value><![CDATA[10]]></value>",
		"<timestamp><![CDATA[2026-05-22T10:00:00Z]]></timestamp>",
		"<value><![CDATA[12]]></value>",
		"<log_module>",
		"<source><![CDATA[application.log]]></source>",
		"<data>&#34;Service restarted&#34;</data>",
		"<encoding>utf-8</encoding>",
	}

	for _, check := range checks {
		if !strings.Contains(xmlString, check) {
			t.Fatalf("expected XML to contain %q, got:\n%s", check, xmlString)
		}
	}

	if strings.Contains(xmlString, "<data><![CDATA[0]]></data>") {
		t.Fatalf("did not expect scalar data node for datalist payload, got:\n%s", xmlString)
	}
}
