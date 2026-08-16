#!/usr/bin/env python3
"""
Parity harness (Python side): builds agent XML and discovery JSON output for
each scenario in fixtures.json using the real pandoraPlugintoolsBasic
library, and prints one JSON object per scenario to stdout for compare.py to
diff against the Go harness's output.

Usage: PYTHONPATH=<path to pandoraPlugintoolsBasic> python3 python_harness.py fixtures.json
"""
import io
import json
import sys
from contextlib import redirect_stdout

import pandoraPlugintools as ppt


def to_python_module_dict(m):
    """Translates the shared fixture schema into what pandoraPlugintoolsBasic
    actually expects:
    - 'alert_templates' (matching the Go Config field name) -> 'alert'
      (the real init_module() template key).
    - 'data_list' (a separate field in the fixture and in Go's Config) ->
      Python has no separate datalist field; print_module detects datalist
      mode by checking isinstance(data["value"], list), so the points must
      be assigned directly to the 'value' key instead.
    """
    m = dict(m)
    if "alert_templates" in m:
        m["alert"] = m.pop("alert_templates")
    if "data_list" in m:
        points = m.pop("data_list")
        m["value"] = [
            {"value": p["value"], **({"timestamp": p["timestamp"]} if p.get("timestamp") else {})}
            for p in points
        ]
    return m


def build_agent_xml(scenario):
    modules = [ppt.init_module(to_python_module_dict(m)) for m in scenario.get("modules", [])]
    log_modules = [ppt.init_log_module(m) for m in scenario.get("log_modules", [])]
    image_modules = [ppt.init_module(to_python_module_dict(m)) for m in scenario.get("image_modules", [])]

    agent = ppt.init_agent(scenario["agent"])

    return ppt.print_agent(
        agent=agent,
        modules=modules,
        log_modules=log_modules,
        log_encoding=scenario.get("log_encoding", ""),
        image_modules=image_modules,
    )


def build_discovery_json(scenario):
    disco = scenario.get("discovery", {})

    ppt.set_disco_error_level(disco.get("error_level", 0))
    ppt.set_disco_summary({})
    ppt.set_disco_info_value("")
    ppt.set_disco_monitoring_data([])

    for item in disco.get("summary_values", []):
        ppt.set_disco_summary_value(item["key"], item["value"])

    for item in disco.get("summary_adds", []):
        ppt.add_disco_summary_value(item["key"], item["value"])

    if disco.get("info_set") is not None:
        ppt.set_disco_info_value(disco["info_set"])

    for text in disco.get("info_adds", []):
        ppt.add_disco_info_value(text)

    for entry in disco.get("monitoring_data_adds", []):
        ppt.add_disco_monitoring_data(entry)

    # Exercise the real disco_output() (it prints JSON and calls
    # sys.exit(_ERROR_LEVEL)) rather than re-reading the module's private
    # globals, so the harness proves the actual production code path.
    buf = io.StringIO()
    with redirect_stdout(buf):
        try:
            ppt.disco_output()
        except SystemExit:
            pass

    return json.loads(buf.getvalue())


def main():
    fixtures_path = sys.argv[1]
    with open(fixtures_path) as f:
        fixtures = json.load(f)

    results = []
    for scenario in fixtures["scenarios"]:
        results.append({
            "name": scenario["name"],
            "agent_xml": build_agent_xml(scenario),
            "discovery_json": build_discovery_json(scenario),
        })

    print(json.dumps(results))


if __name__ == "__main__":
    main()
