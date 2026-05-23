package slicer_test

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/example/logslice/internal/slicer"
)

const sampleLog = `2024-03-01T08:00:00Z level=info msg="startup"
2024-03-01T08:01:00Z level=info msg="connected"
2024-03-01T08:02:00Z level=warn msg="slow query"
2024-03-01T08:03:00Z level=error msg="timeout"
2024-03-01T08:04:00Z level=info msg="recovered"
`

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func mustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSliceMiddleRange(t *testing.T) {
	path := writeTempLog(t, sampleLog)

	var buf bytes.Buffer
	err := slicer.Slice(path, slicer.Options{
		From: mustParse("2024-03-01T08:01:00Z"),
		To:   mustParse("2024-03-01T08:03:00Z"),
	}, &buf)
	if err != nil {
		t.Fatalf("Slice error: %v", err)
	}

	result := buf.String()
	if !strings.Contains(result, "connected") {
		t.Errorf("expected 'connected' in output")
	}
	if !strings.Contains(result, "timeout") {
		t.Errorf("expected 'timeout' in output")
	}
	if strings.Contains(result, "startup") {
		t.Errorf("unexpected 'startup' in output")
	}
	if strings.Contains(result, "recovered") {
		t.Errorf("unexpected 'recovered' in output")
	}
}

func TestSliceEmptyRange(t *testing.T) {
	path := writeTempLog(t, sampleLog)

	var buf bytes.Buffer
	err := slicer.Slice(path, slicer.Options{
		From: mustParse("2024-03-01T09:00:00Z"),
		To:   mustParse("2024-03-01T10:00:00Z"),
	}, &buf)
	if err != nil {
		t.Fatalf("Slice error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for out-of-range query, got: %s", buf.String())
	}
}

func TestSliceFileNotFound(t *testing.T) {
	var buf bytes.Buffer
	err := slicer.Slice("/nonexistent/path.log", slicer.Options{
		From: mustParse("2024-03-01T08:00:00Z"),
		To:   mustParse("2024-03-01T09:00:00Z"),
	}, &buf)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
