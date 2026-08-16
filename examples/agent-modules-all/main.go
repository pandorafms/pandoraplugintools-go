package main

import (
	"fmt"
	"log"

	pptagent "github.com/pandorafms/pandoraplugintools-go/pkg/agent"
	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

func main() {
	ag, err := pptagent.New(pptagent.Config{AgentName: "endpoint-01", AgentAlias: "endpoint-01"})
	if err != nil {
		log.Fatal(err)
	}

	cpu, err := pptmodule.New(pptmodule.Config{
		Name:        "CPU usage",
		Type:        "generic_data",
		Value:       "10",
		Description: "CPU utilization percentage",
		Unit:        "%",
	})
	if err != nil {
		log.Fatal(err)
	}

	logModule, err := pptmodule.NewLog(pptmodule.LogConfig{
		Source: "application.log",
		Value:  "Service restarted",
	})
	if err != nil {
		log.Fatal(err)
	}

	screenshot, err := pptmodule.New(pptmodule.Config{
		Name:  "Screenshot",
		Value: "aGVsbG8=", // base64-encoded image data
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := ag.AddModule(cpu); err != nil {
		log.Fatal(err)
	}

	if err := ag.AddLogModule(logModule); err != nil {
		log.Fatal(err)
	}

	if err := ag.AddImageModule(screenshot); err != nil {
		log.Fatal(err)
	}

	// Element order matches Python's print_agent: modules, then log modules,
	// then image modules.
	xmlData, err := ag.ModulesXMLWithOptions(pptagent.XMLOptions{LogEncoding: "utf-8"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(xmlData))
}
