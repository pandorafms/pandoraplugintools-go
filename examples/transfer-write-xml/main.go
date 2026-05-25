package main

import (
	"fmt"
	"os"

	ppttransfer "github.com/pandorafms/pandoraplugintools-go/pkg/transfer"
)

func main() {
	file, err := ppttransfer.WriteXML([]byte("<agent_data/>\n"), "agent-123", os.TempDir())
	if err != nil {
		panic(err)
	}

	fmt.Println(file)
}
