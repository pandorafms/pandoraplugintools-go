# Examples

The repository includes runnable example files under `examples/`.

## Agent functions

- `examples/agent-new/main.go` — `pptagent.New`
- `examples/agent-add-module/main.go` — `(*pptagent.Agent).AddModule`
- `examples/agent-add-log-module/main.go` — `(*pptagent.Agent).AddLogModule`
- `examples/agent-validate/main.go` — `(*pptagent.Agent).Validate`
- `examples/agent-xml/main.go` — `(*pptagent.Agent).XML`
- `examples/agent-xml-with-options/main.go` — `(*pptagent.Agent).XMLWithOptions`
- `examples/agent-xml-loop/main.go` — `(*pptagent.Agent).XML` with modules added in a loop
- `examples/agent-modules-both/main.go` — `(*pptagent.Agent).ModulesXMLWithOptions` with data and log modules
- `examples/agent-add-image-module/main.go` — `(*pptagent.Agent).AddImageModule`
- `examples/agent-modules-all/main.go` — `(*pptagent.Agent).ModulesXMLWithOptions` with data, log, and image modules together

## Discovery functions

- `examples/discovery-new/main.go` — `pptdiscovery.New`
- `examples/discovery-set-error-level/main.go` — `(*pptdiscovery.Discovery).SetErrorLevel`
- `examples/discovery-set-summary/main.go` — `(*pptdiscovery.Discovery).SetSummary` / `SetSummaryValue`
- `examples/discovery-add-summary-value/main.go` — `(*pptdiscovery.Discovery).AddSummaryValue`
- `examples/discovery-info/main.go` — `(*pptdiscovery.Discovery).SetInfo` / `AddInfo`
- `examples/discovery-monitoring-data/main.go` — `(*pptdiscovery.Discovery).SetMonitoringData` / `AddMonitoringData`
- `examples/discovery-output-json/main.go` — `(*pptdiscovery.Discovery).OutputJSON`
- `examples/discovery-output/main.go` — `(*pptdiscovery.Discovery).Output` (full plugin-style flow)

## Module functions

- `examples/module-new/main.go` — `pptmodule.New`
- `examples/module-datalist/main.go` — `pptmodule.New` with `DataList`
- `examples/module-new-log/main.go` — `pptmodule.NewLog`
- `examples/module-xml-loop/main.go` — `pptmodule.XML` for multiple modules built in a loop
- `examples/module-xml-both/main.go` — `pptmodule.XMLWithOptions` with data and log modules
- `examples/module-validate/main.go` — `(*pptmodule.Module).Validate`
- `examples/logmodule-validate/main.go` — `(*pptmodule.LogModule).Validate`
- `examples/module-image/main.go` — `pptmodule.ImageXML`

## Transfer functions

- `examples/transfer-options-validate/main.go` — `(*ppttransfer.Options).Validate`
- `examples/transfer-write-xml/main.go` — `ppttransfer.WriteXML`
- `examples/transfer-send-local/main.go` — `ppttransfer.Send` in local mode (default data_in)
- `examples/transfer-send-local-custom/main.go` — `ppttransfer.Send` in local mode (custom directories)
- `examples/transfer-send-tentacle/main.go` — `ppttransfer.Send` in Tentacle mode
- `examples/transfer-send-tentacle-retry/main.go` — `ppttransfer.Send` with `Options.Retries`

## Utility functions

- `examples/util-generate-md5/main.go` — `pptutil.GenerateMD5`
- `examples/util-get-os/main.go` — `pptutil.GetOS`
- `examples/util-now/main.go` — `pptutil.Now`
- `examples/util-encode-decode-string/main.go` — `pptutil.EncodeString` / `pptutil.DecodeString`
- `examples/util-parse-int/main.go` — `pptutil.ParseInt`
- `examples/util-parse-float/main.go` — `pptutil.ParseFloat`
- `examples/util-parse-str/main.go` — `pptutil.ParseStr`
- `examples/util-parse-bool/main.go` — `pptutil.ParseBool`
- `examples/util-translate-macros/main.go` — `pptutil.TranslateMacros`
- `examples/util-safe-input-output/main.go` — `pptutil.SafeInput` / `pptutil.SafeOutput`
- `examples/util-parse-configuration/main.go` — `pptutil.ParseConfiguration`
- `examples/util-parse-csv-file/main.go` — `pptutil.ParseCSVFile`

## Output functions

- `examples/output-print-stdout/main.go` — `pptoutput.PrintStdout`
- `examples/output-print-stderr/main.go` — `pptoutput.PrintStderr`
- `examples/output-print-debug/main.go` — `pptoutput.PrintDebug` and `pptoutput.SetDebug`

## Notes

- Public examples use the alias convention `pptagent`, `pptmodule`, `pptdiscovery`, `ppttransfer`, `pptutil`, and `pptoutput`.
- Image modules (`AddImageModule`/`ImageXML`) render after regular and log modules, matching Python's `print_agent(image_modules=...)` element order.
- The Tentacle example is designed to compile and illustrate usage even if `tentacle_client` is not installed locally.
- Transfer defaults: staging uses `os.TempDir()` (usually `/tmp`), local mode data directory uses `/var/spool/pandora/data_in`.
- For fragment-only XML, the `both` examples are the canonical reference. Use `nil` for the group you are not rendering: `pptmodule.XML(modules, nil)` for data-only, or `pptmodule.XMLWithOptions(nil, logModules, opts)` for log-only.
