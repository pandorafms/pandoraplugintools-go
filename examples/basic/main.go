package main

import (
	"context"
	"fmt"
	"os"

	pptagent "github.com/pandorafms/pandoraPlugintoolsGo/pkg/agent"
	pptmodule "github.com/pandorafms/pandoraPlugintoolsGo/pkg/module"
	ppttransfer "github.com/pandorafms/pandoraPlugintoolsGo/pkg/transfer"
)

func main() {
	ag, err := pptagent.New(pptagent.Config{
		AgentName:   "agent-123",
		AgentAlias:  "WIN-SERV",
		Description: "Default Windows server",
	})
	if err != nil {
		panic(err)
	}

	cpu, err := pptmodule.New(pptmodule.Config{
		Name:                 "CPU usage",
		Type:                 "generic_data",
		Value:                "10",
		Description:          "Percentage of CPU utilization",
		Unit:                 "%",
		ModuleGroup:          "system",
		MinWarning:           "60",
		MaxCritical:          "95",
		Status:               "normal",
		CustomID:             "cpu-usage",
		CriticalInstructions: "Check sustained CPU pressure",
		AlertTemplates:       []string{"cpu-warning", "cpu-critical"},
	})
	if err != nil {
		panic(err)
	}

	if err := ag.AddModule(cpu); err != nil {
		panic(err)
	}

	xmlData, err := ag.XML()
	if err != nil {
		panic(err)
	}

	stagingDir, err := os.MkdirTemp("", "ppt-go-staging-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(stagingDir)

	inboxDir, err := os.MkdirTemp("", "ppt-go-inbox-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(inboxDir)

	file, err := ppttransfer.WriteXML(xmlData, ag.Config.AgentName, stagingDir)
	if err != nil {
		panic(err)
	}

	fmt.Println(file)

	if err := ppttransfer.Send(context.Background(), file, ppttransfer.Options{
		Mode:    ppttransfer.ModeLocal,
		DataDir: inboxDir,
	}); err != nil {
		panic(err)
	}
}
