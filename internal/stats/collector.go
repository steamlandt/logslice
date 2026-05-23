package stats

import (
	"io"
)

// CountingWriter wraps an io.Writer and forwards byte counts to Stats.
type CountingWriter struct {
	w    io.Writer
	stats *Stats
}

// NewCountingWriter returns a CountingWriter that delegates to w and records
// every successful write in s.
func NewCountingWriter(w io.Writer, s *Stats) *CountingWriter {
	return &CountingWriter{w: w, stats: s}
}

// Write implements io.Writer.
func (cw *CountingWriter) Write(p []byte) (int, error) {
	n, err := cw.w.Write(p)
	if n > 0 {
		cw.stats.RecordWritten(n)
	}
	return n, err
}

// CountingReader wraps an io.Reader and forwards byte counts to Stats.
type CountingReader struct {
	r     io.Reader
	stats *Stats
}

// NewCountingReader returns a CountingReader that delegates to r and records
// every successful read in s.
func NewCountingReader(r io.Reader, s *Stats) *CountingReader {
	return &CountingReader{r: r, stats: s}
}

// Read implements io.Reader.
func (cr *CountingReader) Read(p []byte) (int, error) {
	n, err := cr.r.Read(p)
	if n > 0 {
		cr.stats.BytesRead += int64(n)
	}
	return n, err
}
