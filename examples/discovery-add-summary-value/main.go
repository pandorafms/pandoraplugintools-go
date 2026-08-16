package main

import (
	"fmt"

	pptdiscovery "github.com/pandorafms/pandoraplugintools-go/pkg/discovery"
)

func main() {
	d := pptdiscovery.New()

	if err := d.AddSummaryValue("total agents", 1); err != nil {
		panic(err)
	}
	if err := d.AddSummaryValue("total agents", 1); err != nil {
		panic(err)
	}

	fmt.Println(d.Summary["total agents"])
}
