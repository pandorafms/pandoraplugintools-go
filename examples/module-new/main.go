package main

import (
	"fmt"

	pptmodule "github.com/pandorafms/pandoraPlugintoolsGo/pkg/module"
)

func main() {
	mod, err := pptmodule.New(pptmodule.Config{
		Name:        "CPU usage",
		Type:        "generic_data",
		Value:       "10",
		Description: "CPU utilization percentage",
		Unit:        "%",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(mod.Config.Name)
}
