package main

import (
	pptdiscovery "github.com/pandorafms/pandoraplugintools-go/pkg/discovery"
)

func main() {
	d := pptdiscovery.New()

	if err := d.AddSummaryValue("Total agents", 1); err != nil {
		panic(err)
	}
	if err := d.AddSummaryValue("Total agents", 1); err != nil {
		panic(err)
	}

	d.AddInfo("Discovered 2 nodes\n")
	d.AddMonitoringData(map[string]any{"module": "node-1"})
	d.AddMonitoringData(map[string]any{"module": "node-2"})

	// Output prints the JSON payload and exits the process with ErrorLevel,
	// mirroring the Python plugin's disco_output() call at the end of a run.
	d.Output()
}
