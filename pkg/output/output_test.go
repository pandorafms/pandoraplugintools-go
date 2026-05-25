package output_test

import (
	"bytes"
	"strings"
	"testing"

	pptoutput "github.com/pandorafms/pandoraplugintools-go/pkg/output"
)

func TestLoggerPrintStdout(t *testing.T) {
	var buf bytes.Buffer
	l := pptoutput.NewLogger(pptoutput.WithStdout(&buf))
	l.PrintStdout("hello %s", "world")

	if !strings.Contains(buf.String(), "hello world") {
		t.Fatalf("expected stdout to contain 'hello world', got %q", buf.String())
	}
}

func TestLoggerPrintStderr(t *testing.T) {
	var buf bytes.Buffer
	l := pptoutput.NewLogger(pptoutput.WithStderr(&buf))
	l.PrintStderr("error %s", "occurred")

	if !strings.Contains(buf.String(), "error occurred") {
		t.Fatalf("expected stderr to contain 'error occurred', got %q", buf.String())
	}
}

func TestLoggerPrintDebugDisabled(t *testing.T) {
	var buf bytes.Buffer
	l := pptoutput.NewLogger(pptoutput.WithStderr(&buf), pptoutput.WithDebug(false))
	l.PrintDebug("should not appear")

	if buf.String() != "" {
		t.Fatalf("expected no debug output when disabled, got %q", buf.String())
	}
}

func TestLoggerPrintDebugEnabled(t *testing.T) {
	var buf bytes.Buffer
	l := pptoutput.NewLogger(pptoutput.WithStderr(&buf), pptoutput.WithDebug(true))
	l.PrintDebug("debug info")

	if !strings.Contains(buf.String(), "[DEBUG] debug info") {
		t.Fatalf("expected debug output with prefix, got %q", buf.String())
	}
}

func TestPackageLevelFunctionsUseDefaultLogger(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	l := pptoutput.NewLogger(pptoutput.WithStdout(&stdout), pptoutput.WithStderr(&stderr), pptoutput.WithDebug(true))

	l.PrintStdout("stdout message")
	l.PrintStderr("stderr message")
	l.PrintDebug("debug message")

	if !strings.Contains(stdout.String(), "stdout message") {
		t.Fatalf("expected stdout message, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "stderr message") {
		t.Fatalf("expected stderr message, got %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "[DEBUG] debug message") {
		t.Fatalf("expected debug message in stderr, got %q", stderr.String())
	}
}
