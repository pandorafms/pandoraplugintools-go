package main

import (
	"fmt"

	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

func main() {
	fmt.Printf("%q\n", pptutil.ParseStr(42))
	fmt.Printf("%q\n", pptutil.ParseStr(nil))
	fmt.Printf("%q\n", pptutil.ParseStr("already a string"))
}
