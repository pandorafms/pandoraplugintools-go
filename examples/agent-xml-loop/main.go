package main

import (
	"fmt"
	"log"

	pptagent "github.com/pandorafms/pandoraplugintools-go/pkg/agent"
	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

func main() {
	// Create a new agent
	ag, err := pptagent.New(pptagent.Config{AgentName: "web-01", AgentAlias: "web-01"})
	if err != nil {
		log.Fatal(err)
	}

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
		})
		if err != nil {
			log.Fatal(err)
		}
		// Append the created module to the modules slice
		modules = append(modules, mod)
	}

	// Add the created modules to the agent
	for _, mod := range modules {
		if err := ag.AddModule(mod); err != nil {
			log.Fatal(err)
		}
	}

	// Generate XML for the agent
	xmlData, err := ag.XML()
	if err != nil {
		log.Fatal(err)
	}
	// Print the generated XML
	fmt.Print(string(xmlData))
}
