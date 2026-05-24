// Package grep provides regex-based line scanning for logslice.
//
// It is designed to be composed with the slicer: first slice a time-bounded
// segment from a large log file, then pass the resulting reader to Scan to
// locate lines matching a pattern within that segment.
//
// Basic usage:
//
//	results, err := grep.Scan(reader, grep.Options{
//		Pattern:    `ERROR|WARN`,
//		MaxResults: 100,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, r := range results {
//		fmt.Printf("line %d: %s\n", r.LineNo, r.Line)
//	}
//
// Setting Invert to true returns lines that do NOT match the pattern,
// which is useful for excluding noise (e.g. DEBUG lines) from a segment.
package grep
