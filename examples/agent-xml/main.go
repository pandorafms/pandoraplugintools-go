package main

import (
	"fmt"

	pptagent "github.com/pandorafms/pandoraplugintools-go/pkg/agent"
	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

func main() {
	ag, err := pptagent.New(pptagent.Config{AgentName: "agent-123", AgentAlias: "WIN-SERV"})
	if err != nil {
		panic(err)
	}

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

	if err := ag.AddModule(mod); err != nil {
		panic(err)
	}

	xmlData, err := ag.XML()
	if err != nil {
		panic(err)
	}

	fmt.Print(string(xmlData))
}
