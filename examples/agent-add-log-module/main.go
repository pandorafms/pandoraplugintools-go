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

	logModule, err := pptmodule.NewLog(pptmodule.LogConfig{Source: "application.log", Value: "Service restarted"})
	if err != nil {
		panic(err)
	}

	if err := ag.AddLogModule(logModule); err != nil {
		panic(err)
	}

	fmt.Println(len(ag.LogModules))
}
