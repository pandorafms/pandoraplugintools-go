package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	ppttransfer "github.com/pandorafms/pandoraPlugintoolsGo/pkg/transfer"
)

func main() {
	stagingDir, err := os.MkdirTemp("", "ppt-send-local-staging-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(stagingDir)

	inboxDir, err := os.MkdirTemp("", "ppt-send-local-inbox-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(inboxDir)

	file, err := ppttransfer.WriteXML([]byte("<agent_data/>\n"), "agent-123", stagingDir)
	if err != nil {
		panic(err)
	}

	if err := ppttransfer.Send(context.Background(), file, ppttransfer.Options{
		Mode:    ppttransfer.ModeLocal,
		DataDir: inboxDir,
	}); err != nil {
		panic(err)
	}

	fmt.Println(filepath.Join(inboxDir, filepath.Base(file)))
}
