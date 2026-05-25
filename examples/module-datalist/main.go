package main

import (
	"fmt"

	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

func main() {
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

	fmt.Println(len(mod.Config.DataList))
}
