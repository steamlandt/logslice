package stats_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/yourorg/logslice/internal/stats"
)

func TestCountingWriter(t *testing.T) {
	s := stats.New()
	var buf bytes.Buffer
	cw := stats.NewCountingWriter(&buf, s)

	_, err := io.WriteString(cw, "hello world")
	if err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}
	if s.BytesWritten != 11 {
		t.Fatalf("expected 11 bytes written, got %d", s.BytesWritten)
	}
	if buf.String() != "hello world" {
		t.Fatalf("unexpected buffer content: %q", buf.String())
	}
}

func TestCountingWriterMultiple(t *testing.T) {
	s := stats.New()
	var buf bytes.Buffer
	cw := stats.NewCountingWriter(&buf, s)

	for _, msg := range []string{"foo", "bar", "baz"} {
		io.WriteString(cw, msg)
	}
	if s.BytesWritten != 9 {
		t.Fatalf("expected 9 bytes written, got %d", s.BytesWritten)
	}
}

func TestCountingReader(t *testing.T) {
	s := stats.New()
	src := bytes.NewBufferString("hello logslice")
	cr := stats.NewCountingReader(src, s)

	out, err := io.ReadAll(cr)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	if string(out) != "hello logslice" {
		t.Fatalf("unexpected content: %q", string(out))
	}
	if s.BytesRead != 14 {
		t.Fatalf("expected 14 bytes read, got %d", s.BytesRead)
	}
}

func TestCountingReaderEmpty(t *testing.T) {
	s := stats.New()
	cr := stats.NewCountingReader(bytes.NewBufferString(""), s)
	io.ReadAll(cr)
	if s.BytesRead != 0 {
		t.Fatalf("expected 0 bytes read, got %d", s.BytesRead)
	}
}
