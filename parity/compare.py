#!/usr/bin/env python3
"""
Parity check: runs the Python harness (against pandoraPlugintoolsBasic) and
the Go harness (against pandoraplugintools-go) over the same fixtures.json,
then compares their outputs SEMANTICALLY (not byte-for-byte):

  - agent_xml is parsed into an element tree and compared structurally,
    ignoring indentation/whitespace and treating a missing XML attribute the
    same as that attribute present with an empty string value (the two
    libraries use different "omit empty" strategies, which is a legitimate
    implementation choice, not a behavioral difference the Pandora server
    would ever observe).
  - discovery_json is compared as parsed JSON (dict/list equality already
    ignores key order).

Usage:
  python3 compare.py \
      --fixtures fixtures.json \
      --pandora-plugintools-basic /path/to/pandoraPlugintoolsBasic \
      --pandoraplugintools-go /path/to/pandoraplugintools-go
"""
import argparse
import json
import os
import subprocess
import sys
import xml.etree.ElementTree as ET

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))


def run_python_harness(fixtures_path, basic_path):
    result = subprocess.run(
        ["python3", "python_harness.py", fixtures_path],
        cwd=SCRIPT_DIR,
        env={"PYTHONPATH": basic_path, "PATH": os.environ.get("PATH", "/usr/bin:/bin")},
        capture_output=True,
        text=True,
    )
    if result.returncode != 0:
        print("Python harness failed:", result.stderr, file=sys.stderr)
        sys.exit(1)
    return json.loads(result.stdout)


def run_go_harness(fixtures_path, go_repo_path):
    result = subprocess.run(
        ["go", "run", "./parity/go_harness", fixtures_path],
        cwd=go_repo_path,
        capture_output=True,
        text=True,
    )
    if result.returncode != 0:
        print("Go harness failed:", result.stderr, file=sys.stderr)
        sys.exit(1)
    return json.loads(result.stdout)


def canon_xml(xml_string):
    """Parses xml_string into a canonical, comparison-friendly tuple tree:
    (tag, sorted attrib items, stripped text, [child tuples...]).
    """
    root = ET.fromstring(xml_string)
    return _canon_element(root)


def _canon_element(el):
    attrib = {k: v for k, v in el.attrib.items()}
    text = (el.text or "").strip()
    # An empty leaf element (no text, no attributes, no children) carries no
    # information, same as that element being entirely absent -- every
    # OTHER optional field on both sides already follows an "omit if empty"
    # rule; drop these so the one Python field that doesn't (extra_data,
    # which checks key presence instead of non-emptiness) doesn't produce a
    # spurious mismatch.
    children = [
        _canon_element(c)
        for c in el
        if not _is_empty_leaf(c)
    ]
    return (el.tag, attrib, text, children)


def _is_empty_leaf(el):
    return (
        not el.attrib
        and len(el) == 0
        and (el.text or "").strip() == ""
    )


def diff_attribs(a, b):
    """Compares two attrib dicts treating a missing key on either side the
    same as that key present with an empty string value.
    """
    keys = set(a) | set(b)
    diffs = []
    for k in sorted(keys):
        va = a.get(k, "")
        vb = b.get(k, "")
        if va != vb:
            diffs.append((k, va, vb))
    return diffs


def compare_canon(py_canon, go_canon, path="agent_data"):
    py_tag, py_attrib, py_text, py_children = py_canon
    go_tag, go_attrib, go_text, go_children = go_canon

    diffs = []

    if py_tag != go_tag:
        diffs.append(f"{path}: tag mismatch: python={py_tag!r} go={go_tag!r}")
        return diffs

    attrib_diffs = diff_attribs(py_attrib, go_attrib)
    for k, va, vb in attrib_diffs:
        diffs.append(f"{path}: attribute {k!r} mismatch: python={va!r} go={vb!r}")

    if py_text != go_text:
        diffs.append(f"{path}: text mismatch: python={py_text!r} go={go_text!r}")

    if len(py_children) != len(go_children):
        diffs.append(
            f"{path}: child count mismatch: python={len(py_children)} go={len(go_children)}"
        )
    else:
        for i, (pc, gc) in enumerate(zip(py_children, go_children)):
            diffs.extend(compare_canon(pc, gc, path=f"{path}/{pc[0]}[{i}]"))

    return diffs


def compare_json(py_json, go_json, path="discovery_json"):
    if py_json != go_json:
        return [f"{path}: mismatch: python={py_json!r} go={go_json!r}"]
    return []


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--fixtures", default="fixtures.json")
    parser.add_argument("--pandora-plugintools-basic", required=True)
    parser.add_argument("--pandoraplugintools-go", required=True)
    args = parser.parse_args()

    fixtures_path = os.path.abspath(os.path.join(SCRIPT_DIR, args.fixtures))

    py_results = {r["name"]: r for r in run_python_harness(fixtures_path, args.pandora_plugintools_basic)}
    go_results = {r["name"]: r for r in run_go_harness(fixtures_path, args.pandoraplugintools_go)}

    all_names = sorted(set(py_results) | set(go_results))
    failures = 0

    for name in all_names:
        if name not in py_results:
            print(f"FAIL {name}: missing from Python harness output")
            failures += 1
            continue
        if name not in go_results:
            print(f"FAIL {name}: missing from Go harness output")
            failures += 1
            continue

        py_r = py_results[name]
        go_r = go_results[name]

        diffs = []
        try:
            py_canon = canon_xml(py_r["agent_xml"])
            go_canon = canon_xml(go_r["agent_xml"])
            diffs.extend(compare_canon(py_canon, go_canon))
        except ET.ParseError as e:
            diffs.append(f"agent_xml: XML parse error: {e}")

        diffs.extend(compare_json(py_r["discovery_json"], go_r["discovery_json"]))

        if diffs:
            print(f"FAIL {name}")
            for d in diffs:
                print(f"  - {d}")
            failures += 1
        else:
            print(f"PASS {name}")

    print()
    print(f"{len(all_names) - failures}/{len(all_names)} scenarios matched")

    sys.exit(1 if failures else 0)


if __name__ == "__main__":
    main()
