package main

import (
	"fmt"

	pptdiscovery "github.com/pandorafms/pandoraplugintools-go/pkg/discovery"
)

func main() {
	d := pptdiscovery.New()
	d.SetSummaryValue("total agents", 1)
	d.SetInfo("done")

	// Unlike Output, OutputJSON does not print or exit(ErrorLevel) — it just
	// returns the payload, which is useful for testing or custom handling.
	out, err := d.OutputJSON()
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}
