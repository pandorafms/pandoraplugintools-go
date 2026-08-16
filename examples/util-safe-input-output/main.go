package main

import (
	"fmt"

	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

func main() {
	encoded := pptutil.SafeInput(`Hello "World" & <tag>`)
	fmt.Println(encoded)

	decoded := pptutil.SafeOutput(encoded)
	fmt.Println(decoded)
}
