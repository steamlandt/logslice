// Package tail implements an efficient reverse-scan reader for log files.
//
// Unlike utilities that read a file from the beginning to find its end,
// tail scans backwards through the file in fixed-size chunks so that only
// the minimum number of bytes are read from disk — O(result) rather than
// O(file size).
//
// Typical usage:
//
//	lines, err := tail.Lines("/var/log/app.log", 50)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, l := range lines {
//		fmt.Println(l)
//	}
//
// The returned slice is ordered chronologically (oldest first), matching
// the natural reading order of the original file.
package tail
