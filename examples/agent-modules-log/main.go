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

	logModule, err := pptmodule.NewLog(pptmodule.LogConfig{
		Source: "application.log",
		Value:  "Service restarted",
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := ag.AddLogModule(logModule); err != nil {
		log.Fatal(err)
	}

	xmlData, err := ag.ModulesXMLWithOptions(pptagent.XMLOptions{LogEncoding: "utf-8"})
	if err != nil {
		log.Fatal(err)
	}

	pptoutput.PrintStdout("%s", string(xmlData))
}
