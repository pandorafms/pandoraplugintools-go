package main

import (
	"context"
	"fmt"
	"log"

	pptagent "github.com/pandorafms/pandoraplugintools-go/pkg/agent"
	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
	pptoutput "github.com/pandorafms/pandoraplugintools-go/pkg/output"
	ppttransfer "github.com/pandorafms/pandoraplugintools-go/pkg/transfer"
	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

func main() {
	serverName := "WIN-SERV"

	ag, err := pptagent.New(pptagent.Config{
		AgentName:   pptutil.GenerateMD5(serverName),
		AgentAlias:  serverName,
		Description: "Default Windows server",
		OSName:      pptutil.GetOS(),
		Timestamp:   pptutil.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}

	cpu, err := pptmodule.New(pptmodule.Config{
		Name:  "CPU usage",
		Type:  "generic_data",
		Value: "10",
		Desc:  "Percentage of CPU utilization",
		Unit:  "%",
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := ag.AddModule(cpu); err != nil {
		log.Fatal(err)
	}

	xmlData, err := ag.XML()
	if err != nil {
		log.Fatal(err)
	}

	file, err := ppttransfer.WriteXML(xmlData, ag.Config.AgentName, "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", xmlData)

	pptoutput.PrintStdout("Written: %s", file)

	if err := ppttransfer.Send(context.Background(), file, ppttransfer.Options{
		Mode:    ppttransfer.ModeTentacle,
		Address: "localhost",
		Port:    41121,
	}); err != nil {
		log.Fatal(err)
	}
}
