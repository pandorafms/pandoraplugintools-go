package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	pptoutput "github.com/pandorafms/pandoraplugintools-go/pkg/output"
)

func GenerateMD5(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

func GetOS() string {
	return runtime.GOOS
}

// PandoraTimestampLayout is the timestamp format expected by the Pandora server.
const PandoraTimestampLayout = "2006/01/02 15:04:05"

func Now(timezone ...string) string {
	loc := time.Local
	if len(timezone) > 0 && timezone[0] != "" {
		if tz, err := time.LoadLocation(timezone[0]); err == nil {
			loc = tz
		}
	}
	return time.Now().In(loc).Format(PandoraTimestampLayout)
}

// EncodeString base64-encodes s. Ports general.py's encode_string.
func EncodeString(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// DecodeString base64-decodes s. Ports general.py's decode_string.
func DecodeString(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ParseInt converts v to an int, returning 0 if v is not a recognized
// numeric, boolean, or string type, or if a string value fails to parse.
// Ports general.py's parse_int(var), which swallows any conversion error.
func ParseInt(v any) int {
	switch t := v.(type) {
	case int:
		return t
	case int64:
		return int(t)
	case float64:
		return int(t)
	case bool:
		if t {
			return 1
		}
		return 0
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(t))
		if err != nil {
			return 0
		}
		return n
	default:
		return 0
	}
}

// ParseFloat converts v to a float64, returning 0 if v is not a recognized
// numeric, boolean, or string type, or if a string value fails to parse.
// Ports general.py's parse_float(var), which swallows any conversion error.
func ParseFloat(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case bool:
		if t {
			return 1
		}
		return 0
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(t), 64)
		if err != nil {
			return 0
		}
		return f
	default:
		return 0
	}
}

// ParseStr converts v to its string representation, matching Go's fmt
// formatting rules. Ports general.py's parse_str(var); note Go's numeric
// formatting (e.g. 3.0 -> "3") differs from Python's str() output, since
// the two languages format numbers differently by design.
func ParseStr(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// ParseBool converts v to a bool using Python-style truthiness: zero
// numbers, empty strings, and nil are false; any other value is true.
// Ports general.py's parse_bool(var) — this is NOT "true"/"false" string
// parsing, e.g. ParseBool("false") returns true because the string is
// non-empty, matching Python's bool("false") behavior.
func ParseBool(v any) bool {
	switch t := v.(type) {
	case nil:
		return false
	case bool:
		return t
	case int:
		return t != 0
	case int64:
		return t != 0
	case float64:
		return t != 0
	case string:
		return t != ""
	default:
		return true
	}
}

// ParseConfiguration parses a "key<separator>value" file into a map, skipping
// blank lines and lines starting with "#". defaultValues fills in any key not
// present in the file. Ports general.py's parse_configuration; unlike the
// Python source it does not abort at the first malformed line (Python's
// unguarded tuple-unpack raises and aborts the whole parse on a line missing
// the separator) — each malformed line is skipped independently instead,
// since parsing the rest of a config file is more useful than silently
// dropping it. A file read error is logged via output.PrintStderr and
// parsing proceeds with an empty config, same as the Python source. There is
// no default file path (Python defaults to
// "/etc/pandora/pandora_server.conf"); callers must pass one explicitly.
func ParseConfiguration(file string, separator string, defaultValues map[string]string) map[string]string {
	if separator == "" {
		separator = " "
	}

	config := map[string]string{}

	data, err := os.ReadFile(file)
	if err != nil {
		pptoutput.PrintStderr("%s", err.Error())
	} else {
		for _, line := range strings.Split(string(data), "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				continue
			}

			parts := strings.SplitN(trimmed, separator, 2)
			if len(parts) != 2 {
				continue
			}

			config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	for option, value := range defaultValues {
		opt := strings.TrimSpace(option)
		if _, ok := config[opt]; !ok {
			config[opt] = strings.TrimSpace(value)
		}
	}

	return config
}

// ParseCSVFile parses a CSV file into a slice of field slices, skipping
// blank lines and lines starting with "#". A line producing fewer than
// countParameters fields is dropped, optionally logging it via
// output.PrintStderr when printErrors is true. Ports general.py's
// parse_csv_file. A file read error is logged via output.PrintStderr and an
// empty slice is returned, same as the Python source.
func ParseCSVFile(file string, separator string, countParameters int, printErrors bool) [][]string {
	if separator == "" {
		separator = ";"
	}

	result := [][]string{}

	data, err := os.ReadFile(file)
	if err != nil {
		pptoutput.PrintStderr("%s", err.Error())
		return result
	}

	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		values := strings.Split(trimmed, separator)
		if len(values) >= countParameters {
			result = append(result, values)
		} else if printErrors {
			pptoutput.PrintStderr("Csv line: %s does not match minimum parameter defined: %d", line, countParameters)
		}
	}

	return result
}

// MacroReplacement is an ordered macro-name/value pair for TranslateMacros.
type MacroReplacement struct {
	Name  string
	Value string
}

// TranslateMacros replaces every occurrence of each macro's Name with its
// Value in data, applying replacements in slice order. Ports general.py's
// translate_macros(macro_dic, data); it takes an ordered slice instead of a
// map because Go map iteration order is randomized and replacement order
// matters when macro names overlap as substrings of one another.
func TranslateMacros(macros []MacroReplacement, data string) string {
	for _, m := range macros {
		data = strings.ReplaceAll(data, m.Name, m.Value)
	}
	return data
}
