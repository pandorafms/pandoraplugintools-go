package main

import (
	"fmt"

	pptagent "github.com/pandorafms/pandoraplugintools-go/pkg/agent"
)

func main() {
	ag, err := pptagent.New(pptagent.Config{AgentName: "agent-123", AgentAlias: "WIN-SERV"})
	if err != nil {
		panic(err)
	}

	if err := ag.Validate(); err != nil {
		panic(err)
	}

	fmt.Println("agent is valid")
}
