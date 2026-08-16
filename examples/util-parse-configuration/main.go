package main

import (
	"fmt"
	"os"
	"path/filepath"

	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

func main() {
	dir, err := os.MkdirTemp("", "ppt-parse-configuration")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "server.conf")
	content := "# comment\nhost 10.0.0.1\nport 41121\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		panic(err)
	}

	config := pptutil.ParseConfiguration(path, " ", map[string]string{
		"timeout": "30",
	})

	fmt.Println(config["host"])
	fmt.Println(config["port"])
	fmt.Println(config["timeout"])
}
