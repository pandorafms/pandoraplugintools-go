package main

import (
	"fmt"

	pptdiscovery "github.com/pandorafms/pandoraplugintools-go/pkg/discovery"
)

func main() {
	d := pptdiscovery.New()

	d.SetMonitoringData([]map[string]any{{"module": "stale"}})
	d.AddMonitoringData(map[string]any{"module": "cpu"})
	d.AddMonitoringData(map[string]any{"module": "mem"})

	fmt.Println(d.MonitoringData)
}
