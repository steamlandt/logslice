package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/stats"
)

func TestRecordScanned(t *testing.T) {
	s := stats.New()
	s.RecordScanned(42)
	s.RecordScanned(10)
	if s.LinesScanned != 2 {
		t.Fatalf("expected 2 scanned, got %d", s.LinesScanned)
	}
	if s.BytesRead != 52 {
		t.Fatalf("expected 52 bytes read, got %d", s.BytesRead)
	}
}

func TestRecordMatchedAndFiltered(t *testing.T) {
	s := stats.New()
	s.RecordMatched()
	s.RecordMatched()
	s.RecordFiltered()
	if s.LinesMatched != 2 {
		t.Fatalf("expected 2 matched, got %d", s.LinesMatched)
	}
	if s.LinesFiltered != 1 {
		t.Fatalf("expected 1 filtered, got %d", s.LinesFiltered)
	}
}

func TestRecordWritten(t *testing.T) {
	s := stats.New()
	s.RecordWritten(100)
	s.RecordWritten(200)
	if s.BytesWritten != 300 {
		t.Fatalf("expected 300 bytes written, got %d", s.BytesWritten)
	}
}

func TestStop(t *testing.T) {
	s := stats.New()
	time.Sleep(2 * time.Millisecond)
	s.Stop()
	if s.Duration < time.Millisecond {
		t.Fatalf("expected duration >= 1ms, got %s", s.Duration)
	}
}

func TestPrint(t *testing.T) {
	s := stats.New()
	s.RecordScanned(50)
	s.RecordMatched()
	s.RecordFiltered()
	s.RecordWritten(48)
	s.Stop()

	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()

	for _, want := range []string{
		"lines scanned",
		"lines matched",
		"lines filtered",
		"bytes read",
		"bytes written",
		"duration",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Print output missing %q", want)
		}
	}
}
