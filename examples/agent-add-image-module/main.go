package main

import (
	"fmt"

	pptagent "github.com/pandorafms/pandoraplugintools-go/pkg/agent"
	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

func main() {
	ag, err := pptagent.New(pptagent.Config{AgentName: "srv-image"})
	if err != nil {
		panic(err)
	}

	imgModule, err := pptmodule.New(pptmodule.Config{
		Name:  "Screenshot",
		Value: "aGVsbG8=", // base64-encoded image data
	})
	if err != nil {
		panic(err)
	}

	if err := ag.AddImageModule(imgModule); err != nil {
		panic(err)
	}

	body, err := ag.XML()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
