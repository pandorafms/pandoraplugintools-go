package main

import (
	"log"

	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
	pptoutput "github.com/pandorafms/pandoraplugintools-go/pkg/output"
)

func main() {
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

	logModule, err := pptmodule.NewLog(pptmodule.LogConfig{
		Source: "application.log",
		Value:  "Service restarted",
	})
	if err != nil {
		log.Fatal(err)
	}

	xmlData, err := pptmodule.XMLWithOptions(
		[]pptmodule.Module{cpu},
		[]pptmodule.LogModule{logModule},
		pptmodule.XMLOptions{LogEncoding: "utf-8"},
	)
	if err != nil {
		log.Fatal(err)
	}

	pptoutput.PrintStdout("%s", string(xmlData))
}
