package multifile_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/config"
	"github.com/logslice/logslice/internal/multifile"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "mflog-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, l := range lines {
		fmt.Fprintln(f, l)
	}
	return f.Name()
}

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func baseConfig(t *testing.T) *config.Config {
	t.Helper()
	return &config.Config{
		Pattern: `^(?P<ts>\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)`,
	}
}

func TestRunMultipleFiles(t *testing.T) {
	lines1 := []string{
		"2024-01-01T10:00:00Z file1 first",
		"2024-01-01T10:01:00Z file1 second",
	}
	lines2 := []string{
		"2024-01-01T10:02:00Z file2 first",
		"2024-01-01T10:03:00Z file2 second",
	}
	p1 := writeTempLog(t, lines1)
	p2 := writeTempLog(t, lines2)

	cfg := baseConfig(t)
	proc := multifile.New(cfg, mustTime("2024-01-01T09:59:00Z"), mustTime("2024-01-01T10:04:00Z"))

	var buf bytes.Buffer
	results := proc.Run([]string{p1, p2}, &buf)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error for %s: %v", r.Path, r.Err)
		}
	}
	out := buf.String()
	if !strings.Contains(out, "file1 first") {
		t.Errorf("expected file1 content in output")
	}
	if !strings.Contains(out, "file2 second") {
		t.Errorf("expected file2 content in output")
	}
}

func TestRunMissingFileReturnsError(t *testing.T) {
	cfg := baseConfig(t)
	proc := multifile.New(cfg, mustTime("2024-01-01T09:00:00Z"), mustTime("2024-01-01T11:00:00Z"))

	missing := filepath.Join(t.TempDir(), "no-such-file.log")
	var buf bytes.Buffer
	results := proc.Run([]string{missing}, &buf)

	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if results[0].Err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestRunPartialFailureContinues(t *testing.T) {
	lines := []string{"2024-01-01T10:00:00Z ok line"}
	good := writeTempLog(t, lines)
	missing := filepath.Join(t.TempDir(), "ghost.log")

	cfg := baseConfig(t)
	proc := multifile.New(cfg, mustTime("2024-01-01T09:00:00Z"), mustTime("2024-01-01T11:00:00Z"))

	var buf bytes.Buffer
	results := proc.Run([]string{missing, good}, &buf)

	if len(results) != 2 {
		t.Fatalf("expected 2 results")
	}
	if results[0].Err == nil {
		t.Error("expected error for missing file")
	}
	if results[1].Err != nil {
		t.Errorf("unexpected error for good file: %v", results[1].Err)
	}
	if !strings.Contains(buf.String(), "ok line") {
		t.Error("expected good file output even after partial failure")
	}
}
