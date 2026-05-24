package grep

import (
	"strings"
	"testing"
)

const sampleLog = `2024-01-01T10:00:00Z INFO  service started
2024-01-01T10:00:01Z DEBUG checking config
2024-01-01T10:00:02Z ERROR connection refused
2024-01-01T10:00:03Z INFO  retry attempt 1
2024-01-01T10:00:04Z ERROR timeout reached
2024-01-01T10:00:05Z INFO  service stopped
`

func TestScanBasicMatch(t *testing.T) {
	results, err := Scan(strings.NewReader(sampleLog), Options{Pattern: "ERROR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].LineNo != 3 {
		t.Errorf("expected line 3, got %d", results[0].LineNo)
	}
	if results[1].LineNo != 5 {
		t.Errorf("expected line 5, got %d", results[1].LineNo)
	}
}

func TestScanMaxResults(t *testing.T) {
	results, err := Scan(strings.NewReader(sampleLog), Options{Pattern: "INFO", MaxResults: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestScanInvert(t *testing.T) {
	results, err := Scan(strings.NewReader(sampleLog), Options{Pattern: "INFO", Invert: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// DEBUG(1) + ERROR(2) = 3 lines
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestScanEmptyPattern(t *testing.T) {
	_, err := Scan(strings.NewReader(sampleLog), Options{Pattern: ""})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestScanInvalidPattern(t *testing.T) {
	_, err := Scan(strings.NewReader(sampleLog), Options{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestScanNoMatches(t *testing.T) {
	results, err := Scan(strings.NewReader(sampleLog), Options{Pattern: "CRITICAL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}
