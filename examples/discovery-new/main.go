package main

import (
	"fmt"

	pptdiscovery "github.com/pandorafms/pandoraplugintools-go/pkg/discovery"
)

func main() {
	d := pptdiscovery.New()
	fmt.Println(d.ErrorLevel, len(d.Summary), len(d.MonitoringData))
}
