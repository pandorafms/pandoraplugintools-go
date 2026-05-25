package main

import (
	pptoutput "github.com/pandorafms/pandoraplugintools-go/pkg/output"
)

func main() {
	pptoutput.PrintStdout("Plugin output: %s", "all good")
}
