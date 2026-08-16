package main

import (
	"fmt"

	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

func main() {
	fmt.Println(pptutil.ParseInt("42"))
	fmt.Println(pptutil.ParseInt(3.9))
	fmt.Println(pptutil.ParseInt("not a number"))
}
