package main

import (
	pptoutput "github.com/pandorafms/pandoraplugintools-go/pkg/output"
)

func main() {
	pptoutput.SetDebug(true)
	pptoutput.PrintDebug("Debug message with debug enabled")

	pptoutput.SetDebug(false)
	pptoutput.PrintDebug("This will not appear")
}
