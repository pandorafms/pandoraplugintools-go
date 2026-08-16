// Demonstrates sending the pandoraplugintools-go monitoring payload to the
// Pandora FMS console API v2 "/monitoring" endpoint using net/http (Go's
// standard library — no third-party HTTP client dependency needed here,
// unlike the Python example this mirrors).
//
// Safe to run without a real server: unless PANDORA_MONITORING_URL is set,
// it only prints the method, URL, headers, and body that would be sent.
package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

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

	url := os.Getenv("PANDORA_MONITORING_URL")
	if url == "" {
		url = "https://<console_host>/pandora_console/api/v2/monitoring"
	}
	apiToken := os.Getenv("PANDORA_API_TOKEN")
	if apiToken == "" {
		apiToken = "<api_token>"
	}

	if os.Getenv("PANDORA_MONITORING_URL") != "" {
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(payload))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Authorization", "Bearer "+apiToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(resp.StatusCode, string(body))
		return
	}

	fmt.Println("Dry run (set PANDORA_MONITORING_URL to actually send this payload):")
	fmt.Println("POST", url)
	fmt.Println("Headers: Authorization: Bearer "+apiToken, "| Content-Type: application/json")
	fmt.Println("Body:", payload)
}
