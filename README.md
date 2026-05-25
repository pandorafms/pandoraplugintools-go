# pandoraPlugintoolsGo

`pandoraPlugintoolsGo` is a Go-first toolkit for Pandora FMS plugin development.

The project is being built in layers:

1. a reusable Go library for developers;
2. a CLI built on top of that library for operators and simple command-line workflows.

## Current status

Phase 1 is focused on closing the core library basics before real CLI work:

- agent and module modeling;
- scalar XML payload generation;
- `.data` file writing;
- transfer abstraction;
- richer module metadata such as grouping, thresholds, FF/event controls, status, instructions, and alert templates;
- scalar compatibility aliases matching the Python package (`Data`, `Desc`, `Alert`);
- datalist payload support;
- log module support.

## Repository layout

```text
cmd/ppt/              CLI entrypoint
pkg/agent/            Agent model and orchestration
pkg/module/           Module model and validation
pkg/transfer/         File writing and transport helpers
internal/pandoraxml/  Internal XML encoder
examples/             Dedicated examples for public functions and common flows
```

## Import convention

Public-facing examples use explicit aliases for consistency:

```go
import (
    pptagent "github.com/pandorafms/pandoraPlugintoolsGo/pkg/agent"
    pptmodule "github.com/pandorafms/pandoraPlugintoolsGo/pkg/module"
    ppttransfer "github.com/pandorafms/pandoraPlugintoolsGo/pkg/transfer"
)
```

This keeps examples and consumer code visually uniform across all public components.

## Examples

The repository includes dedicated runnable example files under `examples/`, including:

- `examples/basic/main.go`
- `examples/agent-new/main.go`
- `examples/agent-add-module/main.go`
- `examples/agent-add-log-module/main.go`
- `examples/agent-validate/main.go`
- `examples/agent-xml/main.go`
- `examples/agent-xml-with-options/main.go`
- `examples/module-new/main.go`
- `examples/module-datalist/main.go`
- `examples/module-new-log/main.go`
- `examples/module-validate/main.go`
- `examples/logmodule-validate/main.go`
- `examples/transfer-options-validate/main.go`
- `examples/transfer-write-xml/main.go`
- `examples/transfer-send-local/main.go`
- `examples/transfer-send-tentacle/main.go`

For the full mapping from public functions to example files, see `docs/examples.md`.

