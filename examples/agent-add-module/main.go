package main

import (
	"fmt"

	pptagent "github.com/pandorafms/pandoraPlugintoolsGo/pkg/agent"
	pptmodule "github.com/pandorafms/pandoraPlugintoolsGo/pkg/module"
)

func main() {
	ag, err := pptagent.New(pptagent.Config{AgentName: "agent-123", AgentAlias: "WIN-SERV"})
	if err != nil {
		panic(err)
	}

	mod, err := pptmodule.New(pptmodule.Config{Name: "CPU usage", Type: "generic_data", Value: "10"})
	if err != nil {
		panic(err)
	}

	if err := ag.AddModule(mod); err != nil {
		panic(err)
	}

	fmt.Println(len(ag.Modules))
}
