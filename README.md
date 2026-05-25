# pandoraplugintools-go

`pandoraplugintools-go` is a Go-first toolkit for Pandora FMS plugin development.

The project is being built in layers:

1. a reusable Go library for developers;
2. a CLI built on top of that library for operators and simple command-line workflows.

## Quick example

```go
package main

import (
    "context"
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

    cpu, err := pptmodule.New(pptmodule.Config{
        Name:  "CPU usage",
        Type:  "generic_data",
        Value: "10",
        Desc:  "Percentage of CPU utilization",
        Unit:  "%",
    })
    if err != nil {
        log.Fatal(err)
    }

    if err := ag.AddModule(cpu); err != nil {
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

    pptoutput.PrintStdout("Written: %s", file)

    if err := ppttransfer.Send(context.Background(), file, ppttransfer.Options{
        Mode:    ppttransfer.ModeTentacle,
        Address: "127.0.0.1",
        Port:    41121,
    }); err != nil {
        log.Fatal(err)
    }
}
```

## CLI usage

```sh
# Run a plugin with local transfer (default data_in: /var/spool/pandora/data_in)
ppt run --agent WIN-SERV

# Run with tentacle transfer
ppt run --agent WIN-SERV --transfer tentacle --address 192.168.1.20 --port 41121

# Enable debug output
ppt run --agent WIN-SERV --debug
```

## Defaults

| Setting | Default |
|---------|---------|
| Staging directory | `os.TempDir()` (usually `/tmp`) |
| Local data directory | `/var/spool/pandora/data_in` |
| Tentacle binary | `tentacle_client` |
| Tentacle port | `41121` |

## Repository layout

```text
cmd/ppt/              CLI entrypoint
pkg/agent/            Agent model and orchestration
pkg/module/           Module model and validation
pkg/transfer/         File writing and transport helpers
pkg/util/             General-purpose helpers (MD5, OS detection, timestamp)
pkg/output/           Output and logging helpers
internal/pandoraxml/  Internal XML encoder
examples/             Dedicated examples for public functions and common flows
```

## Import convention

Public-facing examples use explicit aliases for consistency:

```go
import (
    pptagent "github.com/pandorafms/pandoraplugintools-go/pkg/agent"
    pptmodule "github.com/pandorafms/pandoraplugintools-go/pkg/module"
    ppttransfer "github.com/pandorafms/pandoraplugintools-go/pkg/transfer"
    pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
    pptoutput "github.com/pandorafms/pandoraplugintools-go/pkg/output"
)
```

This keeps examples and consumer code visually uniform across all public components.

## Examples

The repository includes dedicated runnable example files under `examples/`.

For the full mapping from public functions to example files, see `docs/examples.md`.
