package module_test

import (
	"encoding/xml"
	"strings"
	"testing"

	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

func TestNewAppliesDefaults(t *testing.T) {
	mod, err := pptmodule.New(pptmodule.Config{Name: "CPU usage"})
	if err != nil {
		t.Fatalf("expected module to be created, got error: %v", err)
	}

	if mod.Config.Type != "generic_data_string" {
		t.Fatalf("expected default type generic_data_string, got %q", mod.Config.Type)
	}

	if mod.Config.Value != "0" {
		t.Fatalf("expected default value 0, got %q", mod.Config.Value)
	}
}

func TestNewNormalizesCompatibilityAliases(t *testing.T) {
	mod, err := pptmodule.New(pptmodule.Config{
		Name:  "CPU usage",
		Value: "10",
		Data:  "42",
		Desc:  "CPU from compatibility alias",
		Alert: []string{"cpu-warning", "cpu-critical"},
	})
	if err != nil {
		t.Fatalf("expected module to be created, got error: %v", err)
	}

	if mod.Config.Value != "42" {
		t.Fatalf("expected Value to be normalized from Data, got %q", mod.Config.Value)
	}

	if mod.Config.Description != "CPU from compatibility alias" {
		t.Fatalf("expected Description to be normalized from Desc, got %q", mod.Config.Description)
	}

	if len(mod.Config.AlertTemplates) != 2 {
		t.Fatalf("expected AlertTemplates to be normalized from Alert, got %v", mod.Config.AlertTemplates)
	}
}

func TestNewAcceptsDataList(t *testing.T) {
	mod, err := pptmodule.New(pptmodule.Config{
		Name: "Process count",
		DataList: []pptmodule.DataPoint{
			{Value: "10", Timestamp: "2026-05-22T10:00:00Z"},
			{Value: "12"},
		},
	})
	if err != nil {
		t.Fatalf("expected datalist module to be created, got error: %v", err)
	}

	if len(mod.Config.DataList) != 2 {
		t.Fatalf("expected 2 data points, got %d", len(mod.Config.DataList))
	}
}

func TestNewRejectsEmptyDataListPointValue(t *testing.T) {
	if _, err := pptmodule.New(pptmodule.Config{
		Name: "Process count",
		DataList: []pptmodule.DataPoint{
			{Timestamp: "2026-05-22T10:00:00Z"},
		},
	}); err == nil {
		t.Fatal("expected error for datalist point without value")
	}
}

func TestNewLogRequiresSource(t *testing.T) {
	if _, err := pptmodule.NewLog(pptmodule.LogConfig{}); err == nil {
		t.Fatal("expected error for missing log module source")
	}
}

func TestNewRequiresName(t *testing.T) {
	if _, err := pptmodule.New(pptmodule.Config{}); err == nil {
		t.Fatal("expected error for missing module name")
	}
}

func TestXMLRendersModulesWithoutAgentWrapper(t *testing.T) {
	mod, err := pptmodule.New(pptmodule.Config{
		Name:        "CPU usage",
		Type:        "generic_data",
		Value:       "10",
		Description: "CPU utilization percentage",
	})
	if err != nil {
		t.Fatalf("expected module to be created, got error: %v", err)
	}

	logModule, err := pptmodule.NewLog(pptmodule.LogConfig{
		Source: "application.log",
		Value:  "Service restarted",
	})
	if err != nil {
		t.Fatalf("expected log module to be created, got error: %v", err)
	}

	xmlData, err := pptmodule.XMLWithOptions(
		[]pptmodule.Module{mod},
		[]pptmodule.LogModule{logModule},
		pptmodule.XMLOptions{LogEncoding: "utf-8"},
	)
	if err != nil {
		t.Fatalf("expected module XML to be generated, got error: %v", err)
	}

	xmlString := string(xmlData)
	checks := []string{
		"<module>",
		"<![CDATA[CPU usage]]>",
		"<data><![CDATA[10]]></data>",
		"<log_module>",
		"<source><![CDATA[application.log]]></source>",
		"<encoding>utf-8</encoding>",
	}

	for _, check := range checks {
		if !strings.Contains(xmlString, check) {
			t.Fatalf("expected XML to contain %q, got:\n%s", check, xmlString)
		}
	}

	if strings.Contains(xmlString, "<agent_data") {
		t.Fatalf("did not expect agent wrapper, got:\n%s", xmlString)
	}

	if strings.Contains(xmlString, xml.Header) {
		t.Fatalf("did not expect XML header, got:\n%s", xmlString)
	}
}

func TestXMLRejectsInvalidModules(t *testing.T) {
	_, err := pptmodule.XML([]pptmodule.Module{{}}, nil)
	if err == nil {
		t.Fatal("expected error for invalid module XML payload")
	}
}

func TestImageXMLPrefixesDataURI(t *testing.T) {
	mod, err := pptmodule.New(pptmodule.Config{
		Name:  "Screenshot",
		Value: "aGVsbG8=",
	})
	if err != nil {
		t.Fatalf("expected module to be created, got error: %v", err)
	}

	body, err := pptmodule.ImageXML(mod)
	if err != nil {
		t.Fatalf("expected image XML to be generated, got error: %v", err)
	}

	xmlString := string(body)

	if !strings.Contains(xmlString, "<data><![CDATA[data:image/png;base64,aGVsbG8=]]></data>") {
		t.Fatalf("expected data URI prefix in image XML, got:\n%s", xmlString)
	}
}

func TestImageXMLDoesNotMutateOriginalModule(t *testing.T) {
	mod, err := pptmodule.New(pptmodule.Config{
		Name:  "Screenshot",
		Value: "aGVsbG8=",
	})
	if err != nil {
		t.Fatalf("expected module to be created, got error: %v", err)
	}

	if _, err := pptmodule.ImageXML(mod); err != nil {
		t.Fatalf("expected image XML to be generated, got error: %v", err)
	}

	if mod.Config.Value != "aGVsbG8=" {
		t.Fatalf("expected original module Value to be unchanged, got %q", mod.Config.Value)
	}
}

func TestImageXMLRejectsInvalidModule(t *testing.T) {
	_, err := pptmodule.ImageXML(pptmodule.Module{})
	if err == nil {
		t.Fatal("expected error for invalid image module")
	}
}
