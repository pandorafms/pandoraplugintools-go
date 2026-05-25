package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	ppttransfer "github.com/pandorafms/pandoraplugintools-go/pkg/transfer"
)

func main() {
	// Custom staging directory (e.g. for testing or non-standard layouts).
	stagingDir, err := os.MkdirTemp("", "ppt-custom-staging-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(stagingDir)

	// Custom data directory (e.g. for testing or non-standard Pandora installs).
	dataDir, err := os.MkdirTemp("", "ppt-custom-datadir-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dataDir)

	file, err := ppttransfer.WriteXML([]byte("<agent_data/>\n"), "agent-123", stagingDir)
	if err != nil {
		panic(err)
	}

	if err := ppttransfer.Send(context.Background(), file, ppttransfer.Options{
		Mode:    ppttransfer.ModeLocal,
		DataDir: dataDir,
	}); err != nil {
		panic(err)
	}

	fmt.Println(filepath.Join(dataDir, filepath.Base(file)))
}
