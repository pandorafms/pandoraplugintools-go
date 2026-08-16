package main

import (
	"fmt"

	pptdiscovery "github.com/pandorafms/pandoraplugintools-go/pkg/discovery"
)

func main() {
	d := pptdiscovery.New()

	d.SetInfo("Discovered nodes:\n")
	d.AddInfo("- node-1\n")
	d.AddInfo("- node-2\n")

	fmt.Print(d.Info)
}
