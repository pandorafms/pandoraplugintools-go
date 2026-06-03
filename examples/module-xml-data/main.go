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

	xmlData, err := pptmodule.XML([]pptmodule.Module{cpu}, nil)
	if err != nil {
		log.Fatal(err)
	}

	pptoutput.PrintStdout("%s", string(xmlData))
}
