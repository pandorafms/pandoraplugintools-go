package monitoring_test

import (
	"encoding/json"
	"testing"

	pptmonitoring "github.com/pandorafms/pandoraplugintools-go/pkg/monitoring"
)

func TestNewIsEmpty(t *testing.T) {
	m := pptmonitoring.New()

	out, err := m.PayloadJSON()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != "[]" {
		t.Fatalf("expected empty JSON array, got %q", out)
	}
}

func TestAddItemAppends(t *testing.T) {
	m := pptmonitoring.New()

	m.AddItem(
		map[string]any{"agent_name": "nginx-01", "os": "Linux"},
		[]map[string]any{{"name": "nginx_status", "type": "generic_proc", "data": 1}},
	)
	m.AddItem(
		map[string]any{"agent_name": "nginx-02"},
		[]map[string]any{{"name": "nginx_status", "type": "generic_proc", "data": 0}},
	)

	if len(m.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(m.Items))
	}
}

func TestSetItemsReplaces(t *testing.T) {
	m := pptmonitoring.New()
	m.AddItem(map[string]any{"agent_name": "stale"}, nil)

	m.SetItems([]pptmonitoring.Item{
		{AgentData: map[string]any{"agent_name": "fresh"}},
	})

	if len(m.Items) != 1 || m.Items[0].AgentData["agent_name"] != "fresh" {
		t.Fatalf("expected items to be replaced, got %v", m.Items)
	}
}

func TestPayloadJSONShape(t *testing.T) {
	m := pptmonitoring.New()
	m.AddItem(
		map[string]any{"agent_name": "nginx-01", "os": "Linux", "interval": 300},
		[]map[string]any{
			{"name": "nginx_status", "type": "generic_proc", "data": 1},
			{"name": "active_connections", "type": "generic_data", "data": 42},
		},
	)

	out, err := m.PayloadJSON()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed []map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if len(parsed) != 1 {
		t.Fatalf("expected 1 item, got %d", len(parsed))
	}

	agentData, ok := parsed[0]["agent_data"].(map[string]any)
	if !ok {
		t.Fatalf("expected agent_data object, got %v", parsed[0]["agent_data"])
	}
	if agentData["agent_name"] != "nginx-01" {
		t.Fatalf("expected agent_name nginx-01, got %v", agentData["agent_name"])
	}

	moduleData, ok := parsed[0]["module_data"].([]any)
	if !ok || len(moduleData) != 2 {
		t.Fatalf("expected 2 module_data entries, got %v", parsed[0]["module_data"])
	}
}
