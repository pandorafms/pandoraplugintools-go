package main

import (
	"fmt"

	pptagent "github.com/pandorafms/pandoraplugintools-go/pkg/agent"
	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

func main() {
	ag, err := pptagent.New(pptagent.Config{AgentName: "agent-logs", AgentAlias: "agent-logs"})
	if err != nil {
		panic(err)
	}

	mod, err := pptmodule.New(pptmodule.Config{
		Name: "Process count",
		Type: "generic_data",
		DataList: []pptmodule.DataPoint{
			{Value: "10", Timestamp: "2026/05/22 10:00:00"},
			{Value: "12"},
		},
	})
	if err != nil {
		panic(err)
	}

	logModule, err := pptmodule.NewLog(pptmodule.LogConfig{Source: "application.log", Value: "Service restarted"})
	if err != nil {
		panic(err)
	}

	if err := ag.AddModule(mod); err != nil {
		panic(err)
	}

	if err := ag.AddLogModule(logModule); err != nil {
		panic(err)
	}

	xmlData, err := ag.XMLWithOptions(pptagent.XMLOptions{LogEncoding: "utf-8"})
	if err != nil {
		panic(err)
	}

	fmt.Print(string(xmlData))
}
