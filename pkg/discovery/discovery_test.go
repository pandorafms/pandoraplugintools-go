package discovery_test

import (
	"testing"

	pptdiscovery "github.com/pandorafms/pandoraplugintools-go/pkg/discovery"
)

func TestSetSummaryReplacesData(t *testing.T) {
	d := pptdiscovery.New()
	d.SetSummaryValue("stale", 1)

	d.SetSummary(map[string]any{"total agents": 3})

	if len(d.Summary) != 1 {
		t.Fatalf("expected summary to be replaced, got %v", d.Summary)
	}
	if d.Summary["total agents"] != 3 {
		t.Fatalf("expected total agents to be 3, got %v", d.Summary["total agents"])
	}
}

func TestSetSummaryValueSetsKey(t *testing.T) {
	d := pptdiscovery.New()

	d.SetSummaryValue("total agents", 1)

	if d.Summary["total agents"] != 1 {
		t.Fatalf("expected total agents to be 1, got %v", d.Summary["total agents"])
	}
}

func TestSetErrorLevel(t *testing.T) {
	d := pptdiscovery.New()

	d.SetErrorLevel(2)

	if d.ErrorLevel != 2 {
		t.Fatalf("expected error level 2, got %d", d.ErrorLevel)
	}
}

func TestAddInfoAppends(t *testing.T) {
	d := pptdiscovery.New()

	d.SetInfo("line1\n")
	d.AddInfo("line2\n")

	if d.Info != "line1\nline2\n" {
		t.Fatalf("unexpected info: %q", d.Info)
	}
}

func TestSetMonitoringDataReplaces(t *testing.T) {
	d := pptdiscovery.New()
	d.AddMonitoringData(map[string]any{"module": "stale"})

	d.SetMonitoringData([]map[string]any{{"module": "cpu"}})

	if len(d.MonitoringData) != 1 {
		t.Fatalf("expected monitoring data to be replaced, got %v", d.MonitoringData)
	}
}

func TestAddMonitoringDataAppends(t *testing.T) {
	d := pptdiscovery.New()

	d.AddMonitoringData(map[string]any{"module": "cpu"})
	d.AddMonitoringData(map[string]any{"module": "mem"})

	if len(d.MonitoringData) != 2 {
		t.Fatalf("expected 2 monitoring entries, got %d", len(d.MonitoringData))
	}
}
