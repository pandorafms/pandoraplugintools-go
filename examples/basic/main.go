package main

import (
	"context"
	"fmt"
	"log"

	pptagent "github.com/pandorafms/pandoraplugintools-go/pkg/agent"
	pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
	pptoutput "github.com/pandorafms/pandoraplugintools-go/pkg/output"
	ppttransfer "github.com/pandorafms/pandoraplugintools-go/pkg/transfer"
	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

func main() {
	serverName := "WIN-SERV"

	ag, err := pptagent.New(pptagent.Config{
		AgentName:   pptutil.GenerateMD5(serverName),
		AgentAlias:  serverName,
		Description: "Default Windows server",
		OSName:      pptutil.GetOS(),
		Timestamp:   pptutil.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// generic_data: numeric value with warning/critical thresholds, range, and post-processing.
	cpu, err := pptmodule.New(pptmodule.Config{
		Name:                 "CPU usage",
		Type:                 "generic_data",
		Value:                "10",
		Desc:                 "Percentage of CPU utilization",
		Unit:                 "%",
		Interval:             "5",
		Tags:                 "cpu_usage,performance",
		ModuleGroup:          "System",
		Min:                  "0",
		Max:                  "100",
		MinWarning:           "70",
		MaxWarning:           "90",
		MinWarningForced:     "70",
		MaxWarningForced:     "80",
		MinCritical:          "91",
		MaxCritical:          "0",
		MinCriticalForced:    "80",
		MaxCriticalForced:    "100",
		CriticalInverse:      "0",
		WarningInverse:       "0",
		PostProcess:          "1.0",
		Disabled:             "0",
		CriticalInstructions: "Check CPU-intensive processes with top",
		WarningInstructions:  "Monitor load trend",
		UnknownInstructions:  "Verify agent connectivity",
		CustomID:             "cpu-usage",
		ExtraData:            "500",
		AlertTemplates:       []string{"CPU critical alert"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// generic_data_string: string value with pattern-based warning/critical thresholds.
	serviceStatus, err := pptmodule.New(pptmodule.Config{
		Name:               "Service status",
		Type:               "generic_data_string",
		Value:              "running",
		Desc:               "Current state of the web service",
		ModuleGroup:        "System",
		ModuleParent:       "CPU usage",
		ModuleParentUnlink: "0",
		StrWarning:         "degraded",
		StrWarningForced:   "0",
		StrCritical:        "stopped",
		StrCriticalForced:  "0",
		Quiet:              "0",
		CustomID:           "svc-status",
	})
	if err != nil {
		log.Fatal(err)
	}

	// generic_proc: boolean (0/1) with flip-flop event filtering.
	sshProcess, err := pptmodule.New(pptmodule.Config{
		Name:               "SSH process",
		Type:               "generic_proc",
		Value:              "1",
		Desc:               "Whether the sshd process is running",
		ModuleGroup:        "System",
		MinFFEvent:         "2",
		MinFFEventNormal:   "1",
		MinFFEventWarning:  "2",
		MinFFEventCritical: "3",
		FFType:             "0",
		FFTimeout:          "30",
		EachFF:             "0",
		ModuleFFInterval:   "60",
		CronTab:            "*/5 * * * *",
		Disabled:           "0",
	})
	if err != nil {
		log.Fatal(err)
	}

	// generic_data_inc: incremental counter — server computes delta between samples.
	netTraffic, err := pptmodule.New(pptmodule.Config{
		Name:        "Network traffic",
		Type:        "generic_data_inc",
		Value:       "928596884",
		Desc:        "Total bytes transferred (incremental)",
		Unit:        "bytes/sec",
		ModuleGroup: "Networking",
		MinWarning:  "10000000",
		MinCritical: "50000000",
		PostProcess: "1.0",
		Status:      "critical",
		Timestamp:   "",
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := ag.AddModule(cpu); err != nil {
		log.Fatal(err)
	}
	if err := ag.AddModule(serviceStatus); err != nil {
		log.Fatal(err)
	}
	if err := ag.AddModule(sshProcess); err != nil {
		log.Fatal(err)
	}
	if err := ag.AddModule(netTraffic); err != nil {
		log.Fatal(err)
	}

	xmlData, err := ag.XML()
	if err != nil {
		log.Fatal(err)
	}

	file, err := ppttransfer.WriteXML(xmlData, ag.Config.AgentName, "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", xmlData)

	pptoutput.PrintStdout("Written: %s", file)

	if err := ppttransfer.Send(context.Background(), file, ppttransfer.Options{
		Mode:    ppttransfer.ModeTentacle,
		Address: "localhost",
		Port:    41121,
	}); err != nil {
		log.Fatal(err)
	}
}
