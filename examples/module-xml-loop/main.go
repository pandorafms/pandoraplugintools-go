package main

import (
	"fmt"
	"log"

	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

func main() {
	// Slice to hold the created modules
	var modules []pptmodule.Module

	// fake values for the example
	moduleNames := []string{"CPU Usage", "Memory Usage", "Disk Usage"}
	moduleValues := []string{"75", "60", "80"}

	// Loop to create modules based on the names and values
	for i := range moduleNames {
		mod, err := pptmodule.New(pptmodule.Config{
			Name:        moduleNames[i],
			Type:        "generic_data",
			Value:       moduleValues[i],
			Description: "Generated inside a loop",
			Unit:        "%",
		})
		if err != nil {
			log.Fatal(err)
		}
		// Append the created module to the modules slice
		modules = append(modules, mod)
	}
	// Generate XML for the modules
	xmlData, err := pptmodule.XML(modules, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Print the generated XML
	fmt.Print(string(xmlData))
}
