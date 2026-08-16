package main

import (
	"fmt"

	pptdiscovery "github.com/pandorafms/pandoraplugintools-go/pkg/discovery"
)

func main() {
	d := pptdiscovery.New()
	d.SetErrorLevel(1)

	fmt.Println(d.ErrorLevel)
}
