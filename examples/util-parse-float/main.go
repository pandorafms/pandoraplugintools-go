package main

import (
	"fmt"

	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

func main() {
	fmt.Println(pptutil.ParseFloat("2.5"))
	fmt.Println(pptutil.ParseFloat(4))
	fmt.Println(pptutil.ParseFloat("nope"))
}
