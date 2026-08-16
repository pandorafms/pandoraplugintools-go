package main

import (
	"fmt"

	pptmonitoring "github.com/pandorafms/pandoraplugintools-go/pkg/monitoring"
)

func main() {
	m := pptmonitoring.New()

	m.AddItem(
		map[string]any{
			"agent_name": "nginx-01",
			"os":         "Linux",
			"interval":   300,
		},
		[]map[string]any{
			{"name": "nginx_status", "type": "generic_proc", "data": 1},
			{"name": "active_connections", "type": "generic_data", "data": 42},
		},
	)

	payload, err := m.PayloadJSON()
	if err != nil {
		panic(err)
	}

	fmt.Println(payload)
}
