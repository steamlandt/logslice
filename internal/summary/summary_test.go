package summary

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func baseReport() *Report {
	return &Report{
		InputFile: "/var/log/app.log",
		From:      time.Date(2024, 1, 10, 8, 0, 0, 0, time.UTC),
		To:        time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
		Scanned:   1000,
		Matched:   200,
		Filtered:  50,
		Written:   150,
		Elapsed:   123 * time.Millisecond,
	}
}

func TestPrintContainsKeyFields(t *testing.T) {
	r := baseReport()
	var buf bytes.Buffer
	r.Print(&buf)
	out := buf.String()

	for _, want := range []string{
		"/var/log/app.log",
		"2024-01-10T08:00:00Z",
		"2024-01-10T09:00:00Z",
		"1000",
		"200",
		"150",
		"123ms",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Print() output missing %q\ngot:\n%s", want, out)
		}
	}
}

func TestPrintShowsFilteredWhenNonZero(t *testing.T) {
	r := baseReport()
	var buf bytes.Buffer
	r.Print(&buf)
	if !strings.Contains(buf.String(), "Filtered") {
		t.Error("expected Filtered line in output when Filtered > 0")
	}
}

func TestPrintHidesFilteredWhenZero(t *testing.T) {
	r := baseReport()
	r.Filtered = 0
	var buf bytes.Buffer
	r.Print(&buf)
	if strings.Contains(buf.String(), "Filtered") {
		t.Error("expected no Filtered line when Filtered == 0")
	}
}

func TestAddWarning(t *testing.T) {
	r := baseReport()
	r.AddWarning("unparseable line at offset 42")
	var buf bytes.Buffer
	r.Print(&buf)
	out := buf.String()
	if !strings.Contains(out, "unparseable line at offset 42") {
		t.Errorf("expected warning in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Warnings") {
		t.Errorf("expected Warnings header in output, got:\n%s", out)
	}
}

func TestMatchRate(t *testing.T) {
	r := baseReport() // scanned=1000, matched=200
	if got := r.MatchRate(); got != 0.2 {
		t.Errorf("MatchRate() = %v, want 0.2", got)
	}
}

func TestMatchRateZeroScanned(t *testing.T) {
	r := &Report{}
	if got := r.MatchRate(); got != 0 {
		t.Errorf("MatchRate() with 0 scanned = %v, want 0", got)
	}
}
