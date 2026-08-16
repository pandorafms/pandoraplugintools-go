package main

import (
	"fmt"

	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

func main() {
	// Python-style truthiness: any non-empty string is true, even "false".
	fmt.Println(pptutil.ParseBool("false"))
	fmt.Println(pptutil.ParseBool(""))
	fmt.Println(pptutil.ParseBool(0))
	fmt.Println(pptutil.ParseBool(5))
}
