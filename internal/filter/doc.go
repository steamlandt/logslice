// Package filter implements log line filtering for logslice.
//
// A Filter can be configured with a minimum severity level and/or a keyword
// substring. Lines are passed through Match to determine whether they should
// be included in the output slice.
//
// Supported level tokens (case-insensitive): debug, info, warn, warning,
// error, err. Any line whose highest-matched token is below the configured
// minimum level is excluded. When no level or keyword is set every line
// passes through unchanged.
//
// Typical usage:
//
//	f := filter.New("warn", "database")
//	for _, line := range lines {
//		if f.Match(line) {
//			fmt.Println(line)
//		}
//	}
package filter
