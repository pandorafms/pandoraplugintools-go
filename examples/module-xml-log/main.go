package main

import (
	"log"

	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
	pptoutput "github.com/pandorafms/pandoraplugintools-go/pkg/output"
)

func main() {
	logModule, err := pptmodule.NewLog(pptmodule.LogConfig{
		Source: "application.log",
		Value:  "Service restarted",
	})
	if err != nil {
		log.Fatal(err)
	}

	xmlData, err := pptmodule.XMLWithOptions(
		nil,
		[]pptmodule.LogModule{logModule},
		pptmodule.XMLOptions{LogEncoding: "utf-8"},
	)
	if err != nil {
		log.Fatal(err)
	}

	pptoutput.PrintStdout("%s", string(xmlData))
}
