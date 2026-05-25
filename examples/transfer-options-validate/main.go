package main

import (
	"fmt"

	ppttransfer "github.com/pandorafms/pandoraplugintools-go/pkg/transfer"
)

func main() {
	opts := ppttransfer.Options{
		Mode:    ppttransfer.ModeLocal,
		DataDir: "/tmp/pandora-inbox",
	}

	if err := opts.Validate(); err != nil {
		panic(err)
	}

	fmt.Println("transfer options are valid")
}
