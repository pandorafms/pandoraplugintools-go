package main

import (
	"fmt"

	pptdiscovery "github.com/pandorafms/pandoraplugintools-go/pkg/discovery"
)

func main() {
	d := pptdiscovery.New()

	d.SetSummary(map[string]any{"total agents": 3})
	d.SetSummaryValue("deployments agents", 1)

	fmt.Println(d.Summary["total agents"], d.Summary["deployments agents"])
}
