package transfer_test

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	ppttransfer "github.com/pandorafms/pandoraplugintools-go/pkg/transfer"
)

func TestWriteXMLCreatesDataFile(t *testing.T) {
	dir := t.TempDir()

	path, err := ppttransfer.WriteXML([]byte("<agent_data/>"), "agent-123", dir)
	if err != nil {
		t.Fatalf("expected file to be written, got error: %v", err)
	}

	if filepath.Ext(path) != ".data" {
		t.Fatalf("expected .data file, got %q", path)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist, got error: %v", err)
	}
}

func TestWriteXMLCreatesUniqueNames(t *testing.T) {
	dir := t.TempDir()

	first, err := ppttransfer.WriteXML([]byte("<agent_data/>"), "agent-123", dir)
	if err != nil {
		t.Fatalf("expected first file to be written, got error: %v", err)
	}

	second, err := ppttransfer.WriteXML([]byte("<agent_data/>"), "agent-123", dir)
	if err != nil {
		t.Fatalf("expected second file to be written, got error: %v", err)
	}

	if first == second {
		t.Fatalf("expected unique file paths, got %q", first)
	}
}

// writeFakeTentacleClient writes a script that fails failuresBeforeSuccess
// times (tracked via a counter file) before exiting 0, so retry behavior can
// be exercised without a real tentacle_client binary.
func writeFakeTentacleClient(t *testing.T, failuresBeforeSuccess int) string {
	t.Helper()

	dir := t.TempDir()
	counterFile := filepath.Join(dir, "attempts")
	scriptPath := filepath.Join(dir, "fake-tentacle.sh")

	script := "#!/bin/sh\n" +
		"count=$(cat \"" + counterFile + "\" 2>/dev/null || echo 0)\n" +
		"count=$((count + 1))\n" +
		"echo \"$count\" > \"" + counterFile + "\"\n" +
		"if [ \"$count\" -le " + strconv.Itoa(failuresBeforeSuccess) + " ]; then\n" +
		"  echo \"simulated failure $count\" >&2\n" +
		"  exit 1\n" +
		"fi\n" +
		"exit 0\n"

	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		t.Fatalf("failed to write fake tentacle client: %v", err)
	}

	return scriptPath
}

func readAttempts(t *testing.T, scriptPath string) int {
	t.Helper()

	counterFile := filepath.Join(filepath.Dir(scriptPath), "attempts")
	data, err := os.ReadFile(counterFile)
	if err != nil {
		return 0
	}

	n := 0
	for _, c := range data {
		if c < '0' || c > '9' {
			continue
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func TestSendTentacleRetriesUntilSuccess(t *testing.T) {
	dir := t.TempDir()
	sourceFile := filepath.Join(dir, "agent.data")
	if err := os.WriteFile(sourceFile, []byte("payload"), 0o644); err != nil {
		t.Fatalf("expected source file to be created, got error: %v", err)
	}

	script := writeFakeTentacleClient(t, 2)

	err := ppttransfer.Send(context.Background(), sourceFile, ppttransfer.Options{
		Mode:           ppttransfer.ModeTentacle,
		TentacleBinary: script,
		Address:        "127.0.0.1",
		Port:           41121,
		Retries:        2,
	})
	if err != nil {
		t.Fatalf("expected send to succeed after retries, got error: %v", err)
	}

	if attempts := readAttempts(t, script); attempts != 3 {
		t.Fatalf("expected 3 attempts (1 initial + 2 retries), got %d", attempts)
	}
}

func TestSendTentacleFailsWithoutEnoughRetries(t *testing.T) {
	dir := t.TempDir()
	sourceFile := filepath.Join(dir, "agent.data")
	if err := os.WriteFile(sourceFile, []byte("payload"), 0o644); err != nil {
		t.Fatalf("expected source file to be created, got error: %v", err)
	}

	script := writeFakeTentacleClient(t, 5)

	err := ppttransfer.Send(context.Background(), sourceFile, ppttransfer.Options{
		Mode:           ppttransfer.ModeTentacle,
		TentacleBinary: script,
		Address:        "127.0.0.1",
		Port:           41121,
		Retries:        1,
	})
	if err == nil {
		t.Fatalf("expected send to fail when retries are exhausted")
	}

	if attempts := readAttempts(t, script); attempts != 2 {
		t.Fatalf("expected 2 attempts (1 initial + 1 retry), got %d", attempts)
	}
}

func TestSendTentacleDefaultRetriesIsSingleAttempt(t *testing.T) {
	dir := t.TempDir()
	sourceFile := filepath.Join(dir, "agent.data")
	if err := os.WriteFile(sourceFile, []byte("payload"), 0o644); err != nil {
		t.Fatalf("expected source file to be created, got error: %v", err)
	}

	script := writeFakeTentacleClient(t, 1)

	err := ppttransfer.Send(context.Background(), sourceFile, ppttransfer.Options{
		Mode:           ppttransfer.ModeTentacle,
		TentacleBinary: script,
		Address:        "127.0.0.1",
		Port:           41121,
	})
	if err == nil {
		t.Fatalf("expected send to fail on first attempt with default Retries=0")
	}

	if attempts := readAttempts(t, script); attempts != 1 {
		t.Fatalf("expected exactly 1 attempt, got %d", attempts)
	}
}

func TestSendLocalMovesFile(t *testing.T) {
	sourceDir := t.TempDir()
	targetDir := t.TempDir()
	sourceFile := filepath.Join(sourceDir, "agent.data")

	if err := os.WriteFile(sourceFile, []byte("payload"), 0o644); err != nil {
		t.Fatalf("expected source file to be created, got error: %v", err)
	}

	if err := ppttransfer.Send(context.Background(), sourceFile, ppttransfer.Options{
		Mode:    ppttransfer.ModeLocal,
		DataDir: targetDir,
	}); err != nil {
		t.Fatalf("expected local send to succeed, got error: %v", err)
	}

	movedPath := filepath.Join(targetDir, "agent.data")
	if _, err := os.Stat(movedPath); err != nil {
		t.Fatalf("expected moved file to exist, got error: %v", err)
	}

	if _, err := os.Stat(sourceFile); !os.IsNotExist(err) {
		t.Fatalf("expected source file to be gone, got stat error: %v", err)
	}
}
