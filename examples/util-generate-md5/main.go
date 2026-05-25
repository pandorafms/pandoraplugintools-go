package main

import (
	"fmt"

	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

func main() {
	hash := pptutil.GenerateMD5("WIN-SERV")
	fmt.Println(hash)
}
