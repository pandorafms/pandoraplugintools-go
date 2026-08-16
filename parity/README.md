# Parity harness: Python vs Go

Cross-checks pandoraplugintools-go against the real pandoraPlugintoolsBasic
(Python) library by running the same fixtures.json scenarios through both
libraries and comparing the resulting agent XML and discovery JSON
**semantically** (not byte-for-byte): XML is parsed into an element tree and
compared structurally (ignoring indentation/whitespace, and treating a
missing attribute the same as that attribute present with an empty string
value); JSON is compared as parsed data.

Byte-identical output is not the goal or a realistic one: Go's
`encoding/xml.MarshalIndent` and Python's manual string concatenation format
differently by design, and Go's JSON object key order is not guaranteed to
match Python dict insertion order. What matters is that the Pandora server
would receive equivalent data either way.

## Running

```sh
cd parity
python3 compare.py \
  --pandora-plugintools-basic /path/to/pandoraPlugintoolsBasic \
  --pandoraplugintools-go /path/to/pandoraplugintools-go
```

Requires `python3` and `go` on `PATH`. The Python side must be run against a
pandoraPlugintoolsBasic checkout containing the module `value`/`data` key
fix (see PR #7 in that repo) — without it, every scalar module value is
silently zeroed and most scenarios will fail.

## Files

- `fixtures.json` — shared scenario definitions (agent config, modules,
  datalist, log modules, image modules, discovery summary/info/monitoring
  data). Every public configuration parameter is exercised at least once;
  this is deliberately NOT an exhaustive Cartesian product (that would be
  combinatorially explosive), just representative coverage per dimension.
- `python_harness.py` — builds agent XML + discovery JSON for each scenario
  using the real `pandoraPlugintools` package, printing one JSON array to
  stdout. Uses the actual `print_agent`/`disco_output` functions (not
  reimplemented logic) so it exercises the real production code path.
- `go_harness/main.go` — same, using `pandoraplugintools-go`.
- `compare.py` — runs both harnesses and diffs their output.

## Known, accepted differences

- **`extra_data`**: Python's `print_module` checks `"extra_data" in data`
  instead of checking whether it's non-empty, so it always emits an empty
  `<extra_data></extra_data>` after `init_module()`. Go correctly omits it
  when unset. The comparator treats an empty leaf element as equivalent to
  that element being absent.
- **Datalist point timestamps**: Go's `normalizePandoraTimestamp` accepts
  RFC3339 input and converts it to the Pandora format; Python passes
  whatever string it's given through unchanged. Fixtures use the Pandora
  format directly so this Go-only convenience isn't exercised here.
## Bugs found and fixed via this harness

- Go emitted `<min>` before `<max>` for a module's plain threshold fields;
  Python emits `<max>` before `<min>` there (while using min-before-max for
  the warning/critical variants — an inconsistency in the Python source
  itself). Fixed in pandoraplugintools-go PR #11.
- `Agent` had no native way to attach image modules to a full agent XML
  document (unlike Python's `print_agent(image_modules=...)`); `go_harness`
  originally worked around this by splicing `pptmodule.ImageXML()`
  fragments into the document manually. Fixed properly in pandoraplugintools-go
  PR #13 (`Agent.AddImageModule`), and `go_harness` now uses the native API.
