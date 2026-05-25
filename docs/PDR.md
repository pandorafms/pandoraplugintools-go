# pandoraplugintools-go PDR

`pandoraplugintools-go` will be a Go-first toolkit for Pandora FMS plugin development.

The end state is two deliverables built from the same core:

1. A reusable Go library for developers embedding Pandora FMS plugin generation and delivery in their own code.
2. A CLI built on top of that library for operators and users who only need command-line workflows.

This document defines the product goal, non-goals, phased plan, and the architectural direction for the first iterations.

## Quick path

1. Build a small, stable Go library for the most common Pandora plugin workflows.
2. Keep the first versions intentionally narrower than the Python ecosystem.
3. Add a CLI only after the library API proves coherent and reusable.
4. Expand by evidence, not by cloning every Python helper.

## Product goal

Provide a lightweight, idiomatic Go alternative to the basic subset of `pandoraPlugintools`, focused on generating agent/module payloads, producing XML, and delivering those payloads to Pandora FMS with minimal runtime dependencies.

## Problem statement

The current Python tooling is convenient, but it depends on a Python runtime and exposes a dynamic API style that is not ideal for teams standardizing on Go.

A Go implementation should improve:

- deployment simplicity through static binaries;
- concurrency support for discovery and parallel collection tasks;
- API clarity through explicit types and errors;
- reuse across both application code and command-line workflows.

## Final outcome

The project should end with a layered architecture:

- **Core library**: typed Go packages for agent modeling, module modeling, XML generation, transfer, discovery, and shared utilities.
- **CLI**: a thin command-line application that depends on the library and exposes common workflows without duplicating business logic.

## Users

| User | Need |
|------|------|
| Go developers | Embed Pandora FMS payload generation and transfer in services, agents, and tooling |
| Operators / SREs | Generate or send payloads through a CLI without writing Go code |
| Teams migrating from Python | Keep the common workflows while moving to a compiled toolchain |

## Design principles

| Principle | Decision |
|-----------|----------|
| Go-first API | Do not translate the Python API literally; design idiomatic Go types and methods |
| Library first | The CLI must be built on the library, not beside it |
| Small stable core | Start with the high-value subset used most often |
| Explicit behavior | Prefer structs, options, and typed errors over global mutable state |
| Composable layers | Separate build, serialize, persist, and transfer steps |
| Evidence-driven expansion | Add features only after validating real usage or migration demand |

## Non-goals for the first iterations

- Full parity with the entire Python `pandoraPlugintools` ecosystem
- Porting every helper function before validating the core API
- Supporting advanced features with unclear demand in v1
- Designing the CLI as a separate implementation path

## Scope baseline from the Python package

The Python basic package suggests these capability groups:

- general helpers;
- agent creation;
- module creation;
- XML output;
- file transfer;
- parallel execution helpers;
- discovery helpers;
- output/logging helpers.

The Go port should preserve the useful workflow, but not the Python surface shape.

## Proposed Go architecture

A first-pass package layout:

```text
pandoraplugintools-go/
  docs/
  cmd/ppt/
  pkg/agent/
  pkg/module/
  pkg/xml/
  pkg/transfer/
  pkg/discovery/
  pkg/parallel/
  pkg/util/
  internal/
```

### Core concepts

Representative types for the library:

- `Agent`
- `Module`
- `LogModule` (if needed after validation)
- `TransferOptions`
- `DiscoveryReport`
- `Writer` or XML encoder helpers

Representative library operations:

- `NewAgent(config AgentConfig) (*Agent, error)`
- `(*Agent).AddModule(module Module) error`
- `(*Agent).XML() ([]byte, error)`
- `WriteXML(path string, content []byte) error`
- `Send(ctx context.Context, file string, opts TransferOptions) error`

The library should use:

- `context.Context` for long-running or networked operations;
- standard `error` returns;
- options structs instead of long parameter lists;
- minimal hidden global state.

## Phased delivery plan

## Phase 1 — Core library MVP

### Objective
Deliver the smallest useful Go library for the most common Pandora FMS plugin workflows.

### In scope

- agent model and configuration
- module model and configuration
- XML serialization for agent + modules
- file writing helpers
- transfer abstraction with an initial Tentacle-oriented path
- a minimal set of utility helpers that are clearly needed
- examples and tests for the happy path

### Candidate Python-to-Go mappings

| Python concept | Go direction |
|---|---|
| `init_agent` | `NewAgent(AgentConfig)` |
| `print_agent` | `agent.XML()` or XML encoder |
| `init_module` | `NewModule(ModuleConfig)` or struct constructor |
| `write_xml` | `WriteXML(...)` |
| `transfer_xml` / `tentacle_xml` | `Send(...)` with `TransferOptions` |

### Deliverables

- initial Go module
- package skeleton
- typed core API
- unit tests for creation, serialization, and transfer configuration
- examples showing the end-to-end happy path

### Exit criteria

- a developer can generate an agent payload in Go;
- the payload can be written to disk;
- the library API is coherent enough to support a CLI without refactoring fundamentals.

## Phase 2 — Discovery and parallel execution

### Objective
Add the capabilities where Go should provide clear operational value over Python.

### In scope

- discovery report primitives
- worker or goroutine-based concurrency helpers
- timeout/cancellation support through `context.Context`
- aggregation patterns for concurrent collection
- test coverage for concurrent workflows

### Deliverables

- discovery package primitives
- parallel execution helpers
- examples for collection fan-out / fan-in patterns

### Exit criteria

- concurrent plugin collection workflows are possible with a stable API;
- cancellation, timeout, and error propagation are explicit and testable.

## Phase 3 — CLI on top of the library

### Objective
Expose the common workflows through a thin, practical CLI that directly uses the library.

### In scope

- command structure and UX
- config/file input strategy
- payload generation commands
- transfer/send commands
- useful output and error handling

### Candidate command areas

```text
ppt generate
ppt send
ppt validate
ppt discovery
```

### Rules

- CLI commands must call library code, not duplicate it.
- The CLI should remain thin and focused on orchestration, flags, and presentation.

### Deliverables

- `cmd/ppt`
- help text and examples
- end-to-end tests for main command flows

### Exit criteria

- an operator can perform the core workflows without writing Go code;
- the CLI and library stay aligned because they share the same implementation.

## Phase 4 — Compatibility hardening and expansion

### Objective
Close high-value gaps after the core is proven.

### In scope

- migration helpers for teams coming from Python
- additional utility functions with demonstrated demand
- compatibility notes and docs
- packaging and release workflow

### Exit criteria

- the project has a clear migration story;
- feature growth remains controlled and justified.

## Initial success metrics

- the Phase 1 API can express the basic README example from the Python package;
- the library API does not require map-heavy or reflection-heavy escape hatches for normal use;
- the CLI can be implemented as a thin wrapper instead of inventing a second model;
- new features fit the architecture without forcing breaking redesigns.

## Risks and mitigation

| Risk | Mitigation |
|------|------------|
| Trying to replicate Python literally | Review every API against idiomatic Go conventions |
| Scope creep toward full parity too early | Gate new features behind usage evidence and phase boundaries |
| Library and CLI diverge | Make the CLI depend on exported library packages only |
| Transfer behavior becomes too coupled | Keep transport configuration behind clear interfaces/options |
| Weak test coverage | Treat examples and serialization tests as first-class deliverables |

## Decisions for the next step

The next concrete step after this PDR is to define the Phase 1 technical design in more detail:

1. Go module name and repository metadata
2. package structure for the core MVP
3. first public types and constructors
4. XML serialization strategy
5. transfer abstraction boundaries
6. testing strategy for Phase 1

## Checklist

- [x] Final goal defined: reusable library + CLI on top
- [x] Phases defined from MVP to expansion
- [x] Library-first architecture established
- [x] Non-goals documented to prevent premature scope growth
- [x] Phase 1 technical design document
- [x] Initial Go module and package skeleton
- [ ] First implementation tasks
