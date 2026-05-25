package module_test

import (
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
