// Package summary provides a human-readable summary report of a logslice run,
// including time range coverage, line counts, and any warnings encountered.
package summary

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Report holds the collected data for a single logslice operation.
type Report struct {
	InputFile   string
	From        time.Time
	To          time.Time
	Scanned     int64
	Matched     int64
	Filtered    int64
	Written     int64
	Elapsed     time.Duration
	Warnings    []string
}

// AddWarning appends a warning message to the report.
func (r *Report) AddWarning(msg string) {
	r.Warnings = append(r.Warnings, msg)
}

// Print writes a formatted summary to w.
func (r *Report) Print(w io.Writer) {
	fmt.Fprintf(w, "=== logslice summary ===\n")
	fmt.Fprintf(w, "Input   : %s\n", r.InputFile)
	fmt.Fprintf(w, "Range   : %s  →  %s\n",
		r.From.Format(time.RFC3339),
		r.To.Format(time.RFC3339))
	fmt.Fprintf(w, "Scanned : %d lines\n", r.Scanned)
	fmt.Fprintf(w, "Matched : %d lines\n", r.Matched)
	if r.Filtered > 0 {
		fmt.Fprintf(w, "Filtered: %d lines\n", r.Filtered)
	}
	fmt.Fprintf(w, "Written : %d lines\n", r.Written)
	fmt.Fprintf(w, "Elapsed : %s\n", r.Elapsed.Round(time.Millisecond))
	if len(r.Warnings) > 0 {
		fmt.Fprintf(w, "Warnings:\n")
		for _, wn := range r.Warnings {
			fmt.Fprintf(w, "  ! %s\n", wn)
		}
	}
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 24))
}

// MatchRate returns the fraction of scanned lines that fell within the
// requested time range, or 0 if nothing was scanned.
func (r *Report) MatchRate() float64 {
	if r.Scanned == 0 {
		return 0
	}
	return float64(r.Matched) / float64(r.Scanned)
}
