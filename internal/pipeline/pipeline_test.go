package pipeline_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/config"
	"github.com/yourorg/logslice/internal/pipeline"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-*.log")
	if err != nil {
		t.Fatalf("create temp log: %v", err)
	}
	defer f.Close()
	_, _ = f.WriteString(strings.Join(lines, "\n") + "\n")
	return f.Name()
}

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("mustTime: %v", err)
	}
	return v
}

func TestPipelineRun(t *testing.T) {
	lines := []string{
		`2024-01-10T10:00:00Z INFO  startup complete`,
		`2024-01-10T10:01:00Z DEBUG tick`,
		`2024-01-10T10:02:00Z ERROR disk full`,
		`2024-01-10T10:03:00Z INFO  shutdown`,
	}
	input := writeTempLog(t, lines)
	out := filepath.Join(t.TempDir(), "out.log")

	cfg := &config.Config{
		Input:     input,
		Output:    out,
		From:      mustTime(t, "2024-01-10T10:00:00Z"),
		To:        mustTime(t, "2024-01-10T10:03:00Z"),
		ShowStats: false,
	}

	pl, err := pipeline.New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	snap, err := pl.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if snap.Written == 0 {
		t.Error("expected at least one line written")
	}
}

func TestPipelineInvalidInput(t *testing.T) {
	cfg := &config.Config{
		Input:  "/nonexistent/path/file.log",
		Output: "",
		From:   mustTime(t, "2024-01-10T10:00:00Z"),
		To:     mustTime(t, "2024-01-10T10:05:00Z"),
	}

	pl, err := pipeline.New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = pl.Run()
	if err == nil {
		t.Error("expected error for missing input file, got nil")
	}
}

func TestPipelineInvalidPattern(t *testing.T) {
	cfg := &config.Config{
		Input:   "any.log",
		Pattern: "(unclosed",
		From:    mustTime(t, "2024-01-10T10:00:00Z"),
		To:      mustTime(t, "2024-01-10T10:05:00Z"),
	}

	_, err := pipeline.New(cfg)
	if err == nil {
		t.Error("expected error for invalid regex pattern, got nil")
	}
}

// TestPipelineTimeRangeFiltering verifies that lines outside the specified
// time range are excluded from the output.
func TestPipelineTimeRangeFiltering(t *testing.T) {
	lines := []string{
		`2024-01-10T09:59:00Z INFO  before range`,
		`2024-01-10T10:00:00Z INFO  start of range`,
		`2024-01-10T10:01:00Z INFO  inside range`,
		`2024-01-10T10:02:00Z INFO  end of range`,
		`2024-01-10T10:03:00Z INFO  after range`,
	}
	input := writeTempLog(t, lines)
	out := filepath.Join(t.TempDir(), "out.log")

	cfg := &config.Config{
		Input:  input,
		Output: out,
		From:   mustTime(t, "2024-01-10T10:00:00Z"),
		To:     mustTime(t, "2024-01-10T10:02:00Z"),
	}

	pl, err := pipeline.New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	snap, err := pl.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	// Expect only the 3 lines within [10:00, 10:02] to be written.
	const wantWritten = 3
	if snap.Written != wantWritten {
		t.Errorf("Written = %d, want %d", snap.Written, wantWritten)
	}
}
