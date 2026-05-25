package main

import (
	"fmt"

	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

func main() {
	logModule, err := pptmodule.NewLog(pptmodule.LogConfig{
		Source: "application.log",
		Value:  "Service restarted",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(logModule.Config.Source)
}
