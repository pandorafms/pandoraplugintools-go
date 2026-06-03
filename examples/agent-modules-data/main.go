package main

import (
	"log"

	pptagent "github.com/pandorafms/pandoraplugintools-go/pkg/agent"
	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
	pptoutput "github.com/pandorafms/pandoraplugintools-go/pkg/output"
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

	if err := ag.AddModule(cpu); err != nil {
		log.Fatal(err)
	}

	xmlData, err := ag.ModulesXML()
	if err != nil {
		log.Fatal(err)
	}

	pptoutput.PrintStdout("%s", string(xmlData))
}
