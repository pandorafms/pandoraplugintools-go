package main

import (
	"fmt"

	pptagent "github.com/pandorafms/pandoraPlugintoolsGo/pkg/agent"
)

func main() {
	ag, err := pptagent.New(pptagent.Config{
		AgentName:  "agent-123",
		AgentAlias: "WIN-SERV",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(ag.Config.AgentName)
}
