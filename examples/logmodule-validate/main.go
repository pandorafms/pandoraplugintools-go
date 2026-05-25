package main

import (
	"fmt"

	pptmodule "github.com/pandorafms/pandoraPlugintoolsGo/pkg/module"
)

func main() {
	logModule, err := pptmodule.NewLog(pptmodule.LogConfig{Source: "application.log", Value: "Service restarted"})
	if err != nil {
		panic(err)
	}

	if err := logModule.Validate(); err != nil {
		panic(err)
	}

	fmt.Println("log module is valid")
}
