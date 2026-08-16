package util_test

import (
	"testing"
	"time"

	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

func TestGenerateMD5(t *testing.T) {
	result := pptutil.GenerateMD5("WIN-SERV")
	if result == "" {
		t.Fatal("expected non-empty MD5 hash")
	}

	same := pptutil.GenerateMD5("WIN-SERV")
	if result != same {
		t.Fatal("expected same input to produce same hash")
	}

	different := pptutil.GenerateMD5("other")
	if result == different {
		t.Fatal("expected different inputs to produce different hashes")
	}
}

func TestGetOS(t *testing.T) {
	os := pptutil.GetOS()
	if os == "" {
		t.Fatal("expected non-empty OS name")
	}
}

func TestNow(t *testing.T) {
	ts := pptutil.Now()
	if ts == "" {
		t.Fatal("expected non-empty timestamp")
	}
	if _, err := time.Parse(pptutil.PandoraTimestampLayout, ts); err != nil {
		t.Fatalf("expected Pandora timestamp (YYYY/MM/DD HH:MM:SS), got %q: %v", ts, err)
	}
}

func TestNowWithExplicitTimezone(t *testing.T) {
	ts := pptutil.Now("UTC")
	if _, err := time.Parse(pptutil.PandoraTimestampLayout, ts); err != nil {
		t.Fatalf("expected Pandora timestamp (YYYY/MM/DD HH:MM:SS), got %q: %v", ts, err)
	}
}

func TestNowWithInvalidTimezoneUsesLocal(t *testing.T) {
	ts := pptutil.Now("Not/ATimezone")
	if _, err := time.Parse(pptutil.PandoraTimestampLayout, ts); err != nil {
		t.Fatalf("expected Pandora timestamp even with bad timezone, got %q: %v", ts, err)
	}
}

func TestEncodeDecodeStringRoundTrip(t *testing.T) {
	encoded := pptutil.EncodeString("hello world")
	if encoded == "hello world" {
		t.Fatalf("expected encoded string to differ from input")
	}

	decoded, err := pptutil.DecodeString(encoded)
	if err != nil {
		t.Fatalf("unexpected error decoding: %v", err)
	}

	if decoded != "hello world" {
		t.Fatalf("expected round trip to restore original string, got %q", decoded)
	}
}

func TestDecodeStringRejectsInvalidInput(t *testing.T) {
	if _, err := pptutil.DecodeString("not-valid-base64!!"); err == nil {
		t.Fatal("expected error decoding invalid base64")
	}
}

func TestParseInt(t *testing.T) {
	cases := []struct {
		name     string
		input    any
		expected int
	}{
		{"int", 42, 42},
		{"float64", 3.9, 3},
		{"true", true, 1},
		{"false", false, 0},
		{"numeric string", " 42 ", 42},
		{"invalid string", "not a number", 0},
		{"nil", nil, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := pptutil.ParseInt(tc.input); got != tc.expected {
				t.Fatalf("expected %d, got %d", tc.expected, got)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	cases := []struct {
		name     string
		input    any
		expected float64
	}{
		{"float64", 3.5, 3.5},
		{"int", 4, 4.0},
		{"numeric string", "2.5", 2.5},
		{"invalid string", "nope", 0},
		{"nil", nil, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := pptutil.ParseFloat(tc.input); got != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestParseStr(t *testing.T) {
	if got := pptutil.ParseStr("already a string"); got != "already a string" {
		t.Fatalf("expected passthrough for string input, got %q", got)
	}

	if got := pptutil.ParseStr(nil); got != "" {
		t.Fatalf("expected empty string for nil input, got %q", got)
	}

	if got := pptutil.ParseStr(42); got != "42" {
		t.Fatalf("expected \"42\", got %q", got)
	}
}

func TestParseBool(t *testing.T) {
	cases := []struct {
		name     string
		input    any
		expected bool
	}{
		{"nil", nil, false},
		{"zero int", 0, false},
		{"nonzero int", 5, true},
		{"empty string", "", false},
		{"non-empty string is truthy even if it says false", "false", true},
		{"true bool", true, true},
		{"false bool", false, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := pptutil.ParseBool(tc.input); got != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestTranslateMacros(t *testing.T) {
	result := pptutil.TranslateMacros([]pptutil.MacroReplacement{
		{Name: "_host_", Value: "server1"},
		{Name: "_ip_", Value: "10.0.0.1"},
	}, "Host _host_ has IP _ip_")

	if result != "Host server1 has IP 10.0.0.1" {
		t.Fatalf("unexpected result: %q", result)
	}
}

func TestTranslateMacrosAppliesInOrder(t *testing.T) {
	// Later replacements must not re-match text introduced by earlier ones.
	result := pptutil.TranslateMacros([]pptutil.MacroReplacement{
		{Name: "_a_", Value: "_b_"},
		{Name: "_b_", Value: "final"},
	}, "_a_")

	if result != "final" {
		t.Fatalf("expected ordered replacement to cascade, got %q", result)
	}
}

func TestSafeInputEncodesSpecialCharacters(t *testing.T) {
	// Cross-checked against pandoraPlugintoolsBasic's safe_input() output
	// for the same input string.
	got := pptutil.SafeInput("Hello \"World\" & <tag>")
	want := "Hello&#x20;&quot;World&quot;&#x20;&amp;&#x20;&lt;tag&gt;"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestSafeInputLeavesUnmappedCharactersUnchanged(t *testing.T) {
	got := pptutil.SafeInput("plain-text_123")
	if got != "plain-text_123" {
		t.Fatalf("expected unmapped characters to pass through unchanged, got %q", got)
	}
}

func TestSafeOutputDecodesEntitiesBackToCharacters(t *testing.T) {
	got := pptutil.SafeOutput("Hello&#x20;&quot;World&quot;&#x20;&amp;&#x20;&lt;tag&gt;")
	want := "Hello \"World\" & <tag>"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestSafeInputSafeOutputRoundTrip(t *testing.T) {
	sample := "Hello \"World\" & <tag> café niño €" + string(rune(1)) + string(rune(31)) + "\\"

	encoded := pptutil.SafeInput(sample)
	decoded := pptutil.SafeOutput(encoded)

	if decoded != sample {
		t.Fatalf("expected round trip to restore original string, got %q", decoded)
	}
}

func TestSafeInputEmptyString(t *testing.T) {
	if got := pptutil.SafeInput(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestSafeOutputEmptyString(t *testing.T) {
	if got := pptutil.SafeOutput(""); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}
