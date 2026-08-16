package main

import (
	"fmt"

	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

func main() {
	result := pptutil.TranslateMacros([]pptutil.MacroReplacement{
		{Name: "_host_", Value: "server1"},
		{Name: "_ip_", Value: "10.0.0.1"},
	}, "Host _host_ has IP _ip_")

	fmt.Println(result)
}
