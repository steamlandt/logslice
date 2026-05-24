package multifile

import (
	"fmt"
	"io"
)

// PrintSummary writes a human-readable summary of multifile processing
// results to w. It reports per-file line counts and any errors encountered.
func PrintSummary(w io.Writer, results []Result) {
	var totalLines int64
	var errCount int

	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(w, "  [ERR] %s — %v\n", r.Path, r.Err)
			errCount++
			continue
		}
		fmt.Fprintf(w, "  [OK]  %s — %d lines matched\n", r.Path, r.Lines)
		totalLines += r.Lines
	}

	fmt.Fprintf(w, "\nfiles processed: %d, errors: %d, total lines: %d\n",
		len(results), errCount, totalLines)
}
