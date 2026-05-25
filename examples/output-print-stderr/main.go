package main

import (
	pptoutput "github.com/pandorafms/pandoraplugintools-go/pkg/output"
)

func main() {
	pptoutput.PrintStderr("Plugin error: %s", "something failed")
}
