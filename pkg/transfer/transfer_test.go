package transfer_test

import (
	"context"
	"os"
	"path/filepath"
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
