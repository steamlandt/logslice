package grep

import (
	"fmt"
	"io"
)

// PrintOptions controls how results are rendered.
type PrintOptions struct {
	// Numbered prefixes each line with its original line number.
	Numbered bool
	// Highlight wraps the matched portion in ANSI colour codes.
	Highlight bool
	// Pattern is required when Highlight is true.
	Pattern string
}

// Print writes results to w according to opts.
func Print(w io.Writer, results []Result, opts PrintOptions) error {
	for _, r := range results {
		line := r.Line
		if opts.Numbered {
			if _, err := fmt.Fprintf(w, "%6d\t%s\n", r.LineNo, line); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintln(w, line); err != nil {
				return err
			}
		}
	}
	return nil
}

// Summary writes a one-line summary of the scan results to w.
func Summary(w io.Writer, total int, results []Result) error {
	_, err := fmt.Fprintf(w, "scanned %d lines, %d matched\n", total, len(results))
	return err
}
