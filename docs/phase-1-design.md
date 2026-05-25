# Phase 1 Technical Design

Phase 1 will establish the smallest stable Go API that can model a Pandora FMS agent, attach modules, serialize the payload to XML, write it to disk, and prepare delivery through a transfer layer.

## Outcome

At the end of Phase 1, a Go developer should be able to:

1. create an agent with typed configuration;
2. attach one or more modules;
3. serialize the payload to Pandora XML;
4. write the payload to a `.data` file;
5. send or stage the file through a transfer abstraction.

The CLI is intentionally present only as a thin placeholder in this phase. The library remains the primary product.

## Decisions

| Topic | Decision |
|-------|----------|
| Module path | `github.com/pandorafms/pandoraPlugintoolsGo` for the initial scaffold |
| Public layout | Public packages under `pkg/`, XML encoder internals under `internal/` |
| API shape | Typed configs, constructors, methods, and `error` returns |
| XML strategy | Use `encoding/xml` with internal encoder types instead of manual string concatenation |
| Transfer API | A package-level `Send` entrypoint with `TransferOptions` and transport validation |
| File writing | `transfer.WriteXML` writes `.data` payloads and returns the resulting path |
| CLI stance | `cmd/ppt` remains a thin wrapper and must not duplicate library logic |
| Import convention | Public-facing examples use `pptagent`, `pptmodule`, and `ppttransfer` aliases consistently |

## Package layout

```text
cmd/ppt/
pkg/agent/
pkg/module/
pkg/transfer/
internal/pandoraxml/
examples/
```

### Why this layout

- `pkg/agent` owns the aggregate model and public orchestration API.
- `pkg/module` owns reusable module types and validation.
- `pkg/transfer` owns file creation and delivery concerns.
- `internal/pandoraxml` keeps the XML schema details private while allowing the public API to stay stable.
- `cmd/ppt` proves the CLI can sit on top of exported packages later.
- `examples/` provides one-file runnable examples for the public functions and common flows.

## Public API draft

### `pkg/module`

```go
package module

type DataPoint struct {
    Value     string
    Timestamp string
}

type Config struct {
    Name                   string
    Type                   string
    Value                  string
    Data                   string   // compatibility alias -> Value
    DataList               []DataPoint
    Description            string
    Desc                   string   // compatibility alias -> Description
    Unit                   string
    Interval               string
    Tags                   string
    ModuleGroup            string
    ModuleParent           string
    MinWarning             string
    MinWarningForced       string
    MaxWarning             string
    MaxWarningForced       string
    MinCritical            string
    MinCriticalForced      string
    MaxCritical            string
    MaxCriticalForced      string
    StrWarning             string
    StrWarningForced       string
    StrCritical            string
    StrCriticalForced      string
    CriticalInverse        string
    WarningInverse         string
    Min                    string
    Max                    string
    PostProcess            string
    Disabled               string
    MinFFEvent             string
    Status                 string
    Timestamp              string
    CustomID               string
    CriticalInstructions   string
    WarningInstructions    string
    UnknownInstructions    string
    Quiet                  string
    ModuleFFInterval       string
    CronTab                string
    MinFFEventNormal       string
    MinFFEventWarning      string
    MinFFEventCritical     string
    FFType                 string
    FFTimeout              string
    EachFF                 string
    ModuleParentUnlink     string
    ExtraData              string
    AlertTemplates         []string
    Alert                  []string // compatibility alias -> AlertTemplates
}

type Module struct {
    Config Config
}

type LogConfig struct {
    Source string
    Value  string
}

type LogModule struct {
    Config LogConfig
}

func New(cfg Config) (Module, error)
func NewLog(cfg LogConfig) (LogModule, error)
func (m Module) Validate() error
func (m LogModule) Validate() error
```

Design notes:
- Phase 1 starts with the common README workflow, then expands to a first compatibility slice for high-value module metadata.
- The initial expansion includes grouping, parent linkage, warning/critical thresholds, forced thresholds, string/inverse controls, FF/event controls, status/timestamp metadata, instructions, extra data, and alert templates.
- Compatibility aliases from the Python package are accepted for the scalar path: `Data -> Value`, `Desc -> Description`, and `Alert -> AlertTemplates`.
- The scalar path and the structural parity path are both supported now: regular scalar data, datalist payloads, and log modules.
- Later phases can still add discovery output and any remaining lower-use compatibility edges without redesigning the core.

### `pkg/agent`

```go
package agent

type Config struct {
    AgentName       string
    AgentAlias      string
    ParentAgentName string
    Description     string
    Version         string
    OSName          string
    OSVersion       string
    Timestamp       string
    Address         string
    Group           string
    Interval        int
    AgentMode       string
}

type XMLOptions struct {
    LogEncoding string
}

type Agent struct {
    Config     Config
    Modules    []module.Module
    LogModules []module.LogModule
}

func New(cfg Config) (*Agent, error)
func (a *Agent) AddModule(m module.Module) error
func (a *Agent) AddLogModule(m module.LogModule) error
func (a *Agent) XML() ([]byte, error)
func (a *Agent) XMLWithOptions(opts XMLOptions) ([]byte, error)
func (a *Agent) Validate() error
```

Design notes:
- `Agent` is the aggregate root for the Phase 1 happy path.
- Validation should reject missing `AgentName`, invalid scalar/datalist modules, and invalid log modules.
- XML generation stays a method because it is a natural consumer-facing operation.
- `XMLWithOptions` keeps optional log encoding out of the base `XML()` call.

### `pkg/transfer`

```go
package transfer

type Mode string

const (
    ModeTentacle Mode = "tentacle"
    ModeLocal    Mode = "local"
)

type Options struct {
    Mode            Mode
    TentacleBinary  string
    Address         string
    Port            int
    ExtraArgs       []string
    DataDir         string
    RemoveOnSuccess bool
}

func (o Options) Validate() error
func WriteXML(content []byte, agentName string, dir string) (string, error)
func Send(ctx context.Context, file string, opts Options) error
```

Design notes:
- `WriteXML` stays near transfer because the Python package couples file generation and delivery workflows.
- `Send` is transport-oriented and should accept a context.
- Local delivery and Tentacle delivery share the same options object, but validation remains mode-aware.

## XML strategy

The Python implementation builds XML manually with string concatenation. In Go, Phase 1 will use `encoding/xml` through a private package.

### Internal representation

`internal/pandoraxml` will map public agent/module values to private encoder structs.

Benefits:
- escaping and encoding handled by the standard library;
- the public API is decoupled from schema details;
- future support for optional fields remains additive.

### Schema boundary

The public packages should never expose raw encoder structs. They expose domain types and let the internal encoder translate them.

## Validation rules for Phase 1

### Agent
- `AgentName` is required.
- `Interval` defaults to `300` when not set.
- `AgentMode` defaults to `1` when not set.
- `Timestamp` defaults to current UTC time when not set.

### Module
- `Name` is required.
- `Type` defaults to `generic_data` or `generic_data_string` depending on the chosen default policy.
- `Value` defaults to `0` when omitted.

For the initial scaffold we will keep `Type` defaulted to `generic_data_string` to stay close to the Python basic package.

The scalar compatibility slice, datalist payloads, and log-module behavior are now in place. The next major area after closing examples/docs is the real CLI workflow layer.

## Transfer boundaries

Phase 1 supports two delivery modes conceptually:

1. **local**: move the generated `.data` file into a target directory;
2. **tentacle**: invoke the Tentacle client binary with explicit options.

Important boundary decisions:
- transport execution details stay in `pkg/transfer`;
- XML generation does not know anything about delivery;
- the CLI will call `transfer.WriteXML` and `transfer.Send`, not reimplement either concern.

## Testing strategy

Phase 1 tests should cover four layers.

| Layer | Test focus |
|------|------------|
| Module validation | required fields, defaults, invalid configs |
| Agent validation | required agent name, module attachment, defaults |
| XML serialization | stable output shape for the common happy path |
| Transfer helpers | deterministic file naming behavior and option validation |

### Test principles

- Prefer table-driven unit tests.
- Keep XML tests focused on meaningful fields, not formatting trivia.
- Avoid requiring a real Tentacle binary in unit tests.
- Use temporary directories for file-writing tests.

## Example flow to support in Phase 1

```go
import (
    pptagent "github.com/pandorafms/pandoraPlugintoolsGo/pkg/agent"
    pptmodule "github.com/pandorafms/pandoraPlugintoolsGo/pkg/module"
    ppttransfer "github.com/pandorafms/pandoraPlugintoolsGo/pkg/transfer"
)

ag, err := pptagent.New(pptagent.Config{
    AgentName:  "agent-123",
    AgentAlias: "WIN-SERV",
    Description: "Default Windows server",
})
if err != nil {
    return err
}

cpu, err := pptmodule.New(pptmodule.Config{
    Name:        "CPU usage",
    Type:        "generic_data",
    Value:       "10",
    Description: "Percentage of CPU utilization",
    Unit:        "%",
})
if err != nil {
    return err
}

if err := ag.AddModule(cpu); err != nil {
    return err
}

xmlData, err := ag.XML()
if err != nil {
    return err
}

file, err := ppttransfer.WriteXML(xmlData, ag.Config.AgentName, os.TempDir())
if err != nil {
    return err
}

err = ppttransfer.Send(ctx, file, ppttransfer.Options{
    Mode:    ppttransfer.ModeTentacle,
    Address: "192.168.1.20",
    Port:    41121,
})
```

## Deferred items

These are intentionally not part of Phase 1 implementation:

- discovery model details;
- concurrency helpers;
- agent-side convenience CRUD helpers beyond direct slice access;
- advanced or low-use compatibility edge cases beyond the current scalar, datalist, and log-module surface;
- a production-complete CLI command set.

## Exit checklist

- [x] `go.mod` created with the initial module path
- [x] core public packages compile
- [x] internal XML encoder compiles and is hidden from consumers
- [x] minimal example compiles
- [x] unit tests cover validation and serialization basics
- [x] CLI entrypoint builds without duplicating library logic
