// Package monitoring builds the JSON payload accepted by the Pandora FMS
// console API v2 "/monitoring" endpoint. Ports pandoraPlugintoolsBasic's
// monitoring.py.
package monitoring

import "encoding/json"

// Item is one entry of the "/monitoring" API payload array: an agent and
// the modules reported for it. AgentData and ModuleData are plain maps, not
// pptagent.Config/pptmodule.Config, because the monitoring API schema
// differs from the XML/tentacle schema those types model (for example it
// uses "os" instead of "os_name", and module_data.data is a number, not a
// string) — callers build these maps themselves, same as the Python source.
type Item struct {
	AgentData  map[string]any
	ModuleData []map[string]any
}

// Monitoring accumulates the items sent to the "/monitoring" API endpoint.
type Monitoring struct {
	Items []Item
}

// New creates an empty Monitoring.
func New() *Monitoring {
	return &Monitoring{Items: []Item{}}
}

// SetItems replaces the accumulated items with data.
func (m *Monitoring) SetItems(items []Item) {
	m.Items = items
}

// AddItem appends one agent/modules entry to the accumulated items.
func (m *Monitoring) AddItem(agentData map[string]any, moduleData []map[string]any) {
	m.Items = append(m.Items, Item{AgentData: agentData, ModuleData: moduleData})
}

// PayloadJSON builds the JSON array ready to be sent as the request body to
// the "/monitoring" API endpoint.
func (m *Monitoring) PayloadJSON() (string, error) {
	payload := make([]map[string]any, 0, len(m.Items))
	for _, item := range m.Items {
		payload = append(payload, map[string]any{
			"agent_data":  item.AgentData,
			"module_data": item.ModuleData,
		})
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
