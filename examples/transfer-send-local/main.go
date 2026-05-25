package main

import (
	"context"
	"fmt"
	"path/filepath"

	ppttransfer "github.com/pandorafms/pandoraplugintools-go/pkg/transfer"
)

func main() {
	// Using the default staging directory (os.TempDir(), usually /tmp)
	// and the default Pandora data_in directory.
	file, err := ppttransfer.WriteXML([]byte("<agent_data/>\n"), "agent-123", "")
	if err != nil {
		panic(err)
	}

	if err := ppttransfer.Send(context.Background(), file, ppttransfer.Options{
		Mode: ppttransfer.ModeLocal,
	}); err != nil {
		panic(err)
	}

	fmt.Println(filepath.Join("/var/spool/pandora/data_in", filepath.Base(file)))
}
