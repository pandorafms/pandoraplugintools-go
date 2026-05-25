package main

import (
	"context"
	"fmt"
	"os"

	ppttransfer "github.com/pandorafms/pandoraPlugintoolsGo/pkg/transfer"
)

func main() {
	file, err := ppttransfer.WriteXML([]byte("<agent_data/>\n"), "agent-123", os.TempDir())
	if err != nil {
		panic(err)
	}
	defer os.Remove(file)

	err = ppttransfer.Send(context.Background(), file, ppttransfer.Options{
		Mode:           ppttransfer.ModeTentacle,
		TentacleBinary: "tentacle_client",
		Address:        "127.0.0.1",
		Port:           41121,
	})
	if err != nil {
		fmt.Println("tentacle send example result:", err)
		return
	}

	fmt.Println("payload sent")
}
