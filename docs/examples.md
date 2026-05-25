# Examples

The repository includes runnable example files under `examples/`.

## Agent functions

- `examples/agent-new/main.go` — `pptagent.New`
- `examples/agent-add-module/main.go` — `(*pptagent.Agent).AddModule`
- `examples/agent-add-log-module/main.go` — `(*pptagent.Agent).AddLogModule`
- `examples/agent-validate/main.go` — `(*pptagent.Agent).Validate`
- `examples/agent-xml/main.go` — `(*pptagent.Agent).XML`
- `examples/agent-xml-with-options/main.go` — `(*pptagent.Agent).XMLWithOptions`

## Module functions

- `examples/module-new/main.go` — `pptmodule.New`
- `examples/module-datalist/main.go` — `pptmodule.New` with `DataList`
- `examples/module-new-log/main.go` — `pptmodule.NewLog`
- `examples/module-validate/main.go` — `(*pptmodule.Module).Validate`
- `examples/logmodule-validate/main.go` — `(*pptmodule.LogModule).Validate`

## Transfer functions

- `examples/transfer-options-validate/main.go` — `(*ppttransfer.Options).Validate`
- `examples/transfer-write-xml/main.go` — `ppttransfer.WriteXML`
- `examples/transfer-send-local/main.go` — `ppttransfer.Send` in local mode
- `examples/transfer-send-tentacle/main.go` — `ppttransfer.Send` in Tentacle mode

## Notes

- Public examples use the alias convention `pptagent`, `pptmodule`, and `ppttransfer`.
- The Tentacle example is designed to compile and illustrate usage even if `tentacle_client` is not installed locally.
