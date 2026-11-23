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
import multiprocessing
import sys
from contextlib import redirect_stdout

# pandoraPlugintoolsBasic creates a multiprocessing Manager() at import time;
# Python 3.14's forkserver default re-imports this script and dies, so use the
# classic fork start method before importing the package (same workaround as
# the package's own tests).
multiprocessing.set_start_method("fork", force=True)

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

    # Post-transforms run over the complete module list right before the
    # final serialization (ticket #22040 XML path). Tolerant parse: one bad
    # rule is skipped, the good rules still apply.
    pt = scenario.get("post_transforms")
    if pt:
        rules, _ = ppt.parse_transform_blocks_tolerant(pt["rules"])
        ppt.apply_post_transforms(modules, rules)

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

    entries = disco.get("monitoring_data_adds", [])

    # Post-transforms run over the accumulated monitoring items before
    # disco_output() serializes them (ticket #22040 discovery JSON path).
    pt = scenario.get("post_transforms")
    if pt:
        rules, _ = ppt.parse_transform_blocks_tolerant(pt["rules"])
        ppt.apply_post_transforms(entries, rules)

    for entry in entries:
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


def build_post_transform_outputs(scenario):
    """Parser parity and the "/monitoring" JSON path, for scenarios that
    carry post_transforms: the parsed rules (name/match/set, typed values),
    the transformed items built from the raw fixture dicts, and how many
    rule errors the tolerant parser skipped (the error TEXT is per-language,
    so only the count is compared)."""
    pt = scenario.get("post_transforms")
    if not pt:
        return None

    rules, pt_errors = ppt.parse_transform_blocks_tolerant(pt["rules"])

    modules = [to_python_module_dict(dict(m)) for m in scenario.get("modules", [])]
    items = [{"agent_data": dict(scenario["agent"]), "module_data": modules}]
    ppt.apply_post_transforms(items, rules)

    return {
        "parsed_rules": rules,
        "monitoring_items_after": items,
        "post_transform_error_count": len(pt_errors),
    }


def main():
    fixtures_path = sys.argv[1]
    with open(fixtures_path) as f:
        fixtures = json.load(f)

    results = []
    for scenario in fixtures["scenarios"]:
        result = {
            "name": scenario["name"],
            "agent_xml": build_agent_xml(scenario),
            "discovery_json": build_discovery_json(scenario),
        }

        pt_outputs = build_post_transform_outputs(scenario)
        if pt_outputs is not None:
            result.update(pt_outputs)

        results.append(result)

    print(json.dumps(results))


if __name__ == "__main__":
    main()
