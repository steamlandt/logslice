package output

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteLineRaw(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.log")
	w, err := New(Options{Format: FormatRaw, Destination: tmp})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	lines := []string{"line one", "line two", "line three"}
	for i, l := range lines {
		if err := w.WriteLine(i+1, l); err != nil {
			t.Fatalf("WriteLine: %v", err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	got := string(data)
	for _, l := range lines {
		if !strings.Contains(got, l) {
			t.Errorf("expected %q in output", l)
		}
	}
	if strings.Contains(got, "\t") {
		t.Error("raw format should not contain tab separators")
	}
}

func TestWriteLineNumbered(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "numbered.log")
	w, err := New(Options{Format: FormatNumbered, Destination: tmp})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := w.WriteLine(42, "important event"); err != nil {
		t.Fatalf("WriteLine: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	got := string(data)
	if !strings.Contains(got, "42\t") {
		t.Errorf("expected line number prefix, got: %q", got)
	}
	if !strings.Contains(got, "important event") {
		t.Errorf("expected line content, got: %q", got)
	}
}

func TestNewInvalidDestination(t *testing.T) {
	_, err := New(Options{Destination: "/nonexistent/dir/out.log"})
	if err == nil {
		t.Fatal("expected error for invalid destination, got nil")
	}
}

func TestStdoutWriter(t *testing.T) {
	// Ensure a stdout writer (empty destination) constructs without error.
	w, err := New(Options{Format: FormatRaw, Destination: ""})
	if err != nil {
		t.Fatalf("New stdout: %v", err)
	}
	// Close should be a no-op (no owned file).
	if err := w.Close(); err != nil {
		t.Fatalf("Close stdout: %v", err)
	}
}
