// Package discovery ports pandoraPlugintools/discovery.py: it accumulates the
// summary, info, and monitoring data emitted by Pandora FMS discovery plugins
// and renders the final JSON payload consumed by the discovery server.
package discovery

import "fmt"

// Discovery accumulates the state a discovery plugin reports back to the
// Pandora FMS discovery server.
type Discovery struct {
	ErrorLevel     int
	Summary        map[string]any
	Info           string
	MonitoringData []map[string]any
}

// New creates an empty Discovery.
func New() *Discovery {
	return &Discovery{
		Summary:        map[string]any{},
		MonitoringData: []map[string]any{},
	}
}

// SetErrorLevel sets the exit code reported by Output.
func (d *Discovery) SetErrorLevel(value int) {
	d.ErrorLevel = value
}

// SetSummary replaces the summary with the given data.
func (d *Discovery) SetSummary(data map[string]any) {
	d.Summary = data
}

// SetSummaryValue sets a fixed value for a summary key.
func (d *Discovery) SetSummaryValue(key string, value any) {
	d.Summary[key] = value
}

// AddSummaryValue adds value to the existing summary entry at key, or sets it
// if the key is not present yet. Both values must be int, float64, or string;
// mismatched or unsupported types return an error instead of failing silently.
func (d *Discovery) AddSummaryValue(key string, value any) error {
	existing, ok := d.Summary[key]
	if !ok {
		d.SetSummaryValue(key, value)
		return nil
	}

	switch e := existing.(type) {
	case int:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("discovery: cannot add %T to int summary value %q", value, key)
		}
		d.Summary[key] = e + v
	case float64:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("discovery: cannot add %T to float64 summary value %q", value, key)
		}
		d.Summary[key] = e + v
	case string:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("discovery: cannot add %T to string summary value %q", value, key)
		}
		d.Summary[key] = e + v
	default:
		return fmt.Errorf("discovery: unsupported summary value type %T for key %q", existing, key)
	}

	return nil
}

// SetInfo replaces the info string.
func (d *Discovery) SetInfo(value string) {
	d.Info = value
}

// AddInfo appends to the info string.
func (d *Discovery) AddInfo(value string) {
	d.Info += value
}

// SetMonitoringData replaces the monitoring data.
func (d *Discovery) SetMonitoringData(data []map[string]any) {
	d.MonitoringData = data
}

// AddMonitoringData appends an entry to the monitoring data.
func (d *Discovery) AddMonitoringData(data map[string]any) {
	d.MonitoringData = append(d.MonitoringData, data)
}
