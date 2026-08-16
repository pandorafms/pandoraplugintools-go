// Package discovery ports pandoraPlugintools/discovery.py: it accumulates the
// summary, info, and monitoring data emitted by Pandora FMS discovery plugins
// and renders the final JSON payload consumed by the discovery server.
package discovery

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
