package tail

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	for _, l := range lines {
		fmt.Fprintln(f, l)
	}
	return f.Name()
}

func TestLinesBasic(t *testing.T) {
	input := []string{"line1", "line2", "line3", "line4", "line5"}
	path := writeTempLog(t, input)

	got, err := Lines(path, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if got[0] != "line3" || got[1] != "line4" || got[2] != "line5" {
		t.Errorf("unexpected lines: %v", got)
	}
}

func TestLinesMoreThanFile(t *testing.T) {
	input := []string{"alpha", "beta"}
	path := writeTempLog(t, input)

	got, err := Lines(path, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestLinesLargeFile(t *testing.T) {
	var input []string
	for i := 0; i < 500; i++ {
		input = append(input, fmt.Sprintf("2024-01-01T00:00:%02d log entry number %d", i%60, i))
	}
	path := writeTempLog(t, input)

	got, err := Lines(path, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 20 {
		t.Fatalf("expected 20 lines, got %d", len(got))
	}
	if !strings.Contains(got[0], "entry number 480") {
		t.Errorf("first of last-20 should be entry 480, got: %s", got[0])
	}
}

func TestLinesInvalidN(t *testing.T) {
	_, err := Lines("any.log", 0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestLinesFileNotFound(t *testing.T) {
	_, err := Lines(filepath.Join(t.TempDir(), "no-such-file.log"), 5)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
