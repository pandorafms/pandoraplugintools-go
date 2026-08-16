package main

import (
	"fmt"

	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
)

func main() {
	mod, err := pptmodule.New(pptmodule.Config{
		Name:  "Screenshot",
		Value: "aGVsbG8=", // base64-encoded image data
	})
	if err != nil {
		panic(err)
	}

	body, err := pptmodule.ImageXML(mod)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}
