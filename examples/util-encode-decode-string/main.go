package main

import (
	"fmt"

	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

func main() {
	encoded := pptutil.EncodeString("hello world")
	fmt.Println(encoded)

	decoded, err := pptutil.DecodeString(encoded)
	if err != nil {
		panic(err)
	}
	fmt.Println(decoded)
}
