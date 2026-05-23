// Package stats collects and reports metrics about a log slicing operation.
package stats

import (
	"fmt"
	"io"
	"time"
)

// Stats holds counters accumulated during a slice run.
type Stats struct {
	LinesScanned  int64
	LinesMatched  int64
	LinesFiltered int64
	BytesRead     int64
	BytesWritten  int64
	Duration      time.Duration

	start time.Time
}

// New returns a new Stats instance with the clock started.
func New() *Stats {
	return &Stats{start: time.Now()}
}

// RecordScanned increments the scanned-line counter and byte counter.
func (s *Stats) RecordScanned(lineBytes int) {
	s.LinesScanned++
	s.BytesRead += int64(lineBytes)
}

// RecordMatched increments the matched-line counter.
func (s *Stats) RecordMatched() {
	s.LinesMatched++
}

// RecordFiltered increments the filtered-out counter.
func (s *Stats) RecordFiltered() {
	s.LinesFiltered++
}

// RecordWritten increments the written-byte counter.
func (s *Stats) RecordWritten(n int) {
	s.BytesWritten += int64(n)
}

// Stop marks the end of the operation and records elapsed time.
func (s *Stats) Stop() {
	s.Duration = time.Since(s.start)
}

// Print writes a human-readable summary to w.
func (s *Stats) Print(w io.Writer) {
	fmt.Fprintf(w, "lines scanned : %d\n", s.LinesScanned)
	fmt.Fprintf(w, "lines matched : %d\n", s.LinesMatched)
	fmt.Fprintf(w, "lines filtered: %d\n", s.LinesFiltered)
	fmt.Fprintf(w, "bytes read    : %d\n", s.BytesRead)
	fmt.Fprintf(w, "bytes written : %d\n", s.BytesWritten)
	fmt.Fprintf(w, "duration      : %s\n", s.Duration.Round(time.Millisecond))
}
