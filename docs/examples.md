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
- `examples/transfer-send-local/main.go` — `ppttransfer.Send` in local mode (default data_in)
- `examples/transfer-send-local-custom/main.go` — `ppttransfer.Send` in local mode (custom directories)
- `examples/transfer-send-tentacle/main.go` — `ppttransfer.Send` in Tentacle mode

## Utility functions

- `examples/util-generate-md5/main.go` — `pptutil.GenerateMD5`
- `examples/util-get-os/main.go` — `pptutil.GetOS`
- `examples/util-now/main.go` — `pptutil.Now`

## Output functions

- `examples/output-print-stdout/main.go` — `pptoutput.PrintStdout`
- `examples/output-print-stderr/main.go` — `pptoutput.PrintStderr`
- `examples/output-print-debug/main.go` — `pptoutput.PrintDebug` and `pptoutput.SetDebug`

## Notes

- Public examples use the alias convention `pptagent`, `pptmodule`, `ppttransfer`, `pptutil`, and `pptoutput`.
- The Tentacle example is designed to compile and illustrate usage even if `tentacle_client` is not installed locally.
- Transfer defaults: staging uses `os.TempDir()` (usually `/tmp`), local mode data directory uses `/var/spool/pandora/data_in`.
