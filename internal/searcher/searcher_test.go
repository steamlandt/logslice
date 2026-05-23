package searcher_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/example/logslice/internal/searcher"
)

func makeLog(entries []string) *strings.Reader {
	return strings.NewReader(strings.Join(entries, "\n") + "\n")
}

func parseTS(line []byte) (time.Time, error) {
	parts := bytes.SplitN(line, []byte(" "), 2)
	if len(parts) == 0 {
		return time.Time{}, fmt.Errorf("empty line")
	}
	return time.Parse(time.RFC3339, string(parts[0]))
}

var logLines = []string{
	"2024-01-01T10:00:00Z msg=alpha",
	"2024-01-01T10:01:00Z msg=beta",
	"2024-01-01T10:02:00Z msg=gamma",
	"2024-01-01T10:03:00Z msg=delta",
	"2024-01-01T10:04:00Z msg=epsilon",
}

func newSearcher(lines []string) (*searcher.Searcher, int64) {
	content := strings.Join(lines, "\n") + "\n"
	r := strings.NewReader(content)
	return searcher.New(r, int64(len(content)), parseTS), int64(len(content))
}

func TestFindStart(t *testing.T) {
	s, _ := newSearcher(logLines)
	target, _ := time.Parse(time.RFC3339, "2024-01-01T10:02:00Z")

	pos, err := s.FindStart(target)
	if err != nil {
		t.Fatalf("FindStart error: %v", err)
	}
	if pos < 0 {
		t.Fatalf("expected non-negative position, got %d", pos)
	}
}

func TestFindEnd(t *testing.T) {
	s, size := newSearcher(logLines)
	target, _ := time.Parse(time.RFC3339, "2024-01-01T10:03:00Z")

	pos, err := s.FindEnd(target)
	if err != nil {
		t.Fatalf("FindEnd error: %v", err)
	}
	if int64(pos) > size {
		t.Fatalf("end position %d exceeds file size %d", pos, size)
	}
}

func TestFindStartBeforeAll(t *testing.T) {
	s, _ := newSearcher(logLines)
	target, _ := time.Parse(time.RFC3339, "2024-01-01T09:00:00Z")

	pos, err := s.FindStart(target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pos != 0 {
		t.Fatalf("expected position 0, got %d", pos)
	}
}

func TestFindEndAfterAll(t *testing.T) {
	s, size := newSearcher(logLines)
	target, _ := time.Parse(time.RFC3339, "2024-01-01T23:00:00Z")

	pos, err := s.FindEnd(target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if int64(pos) != size {
		t.Fatalf("expected position %d (EOF), got %d", size, pos)
	}
}
