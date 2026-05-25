package main

import (
	"fmt"

	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

func main() {
	mod, err := pptmodule.New(pptmodule.Config{Name: "CPU usage", Type: "generic_data", Value: "10"})
	if err != nil {
		panic(err)
	}

	if err := mod.Validate(); err != nil {
		panic(err)
	}

	fmt.Println("module is valid")
}
