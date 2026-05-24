// Package multifile processes multiple log files in sequence and merges
// their time-bounded output into a single stream.
//
// Typical usage:
//
//	paths := []string{"/var/log/app.log.2", "/var/log/app.log.1", "/var/log/app.log"}
//	p := multifile.New(cfg, from, to)
//	results := p.Run(paths, os.Stdout)
//	for _, r := range results {
//		if r.Err != nil {
//			log.Printf("warn: %v", r.Err)
//		}
//	}
//
// Files are processed in the order provided. Pair this package with
// internal/rotated to discover and sort rotated log segments before
// passing them to the Processor.
package multifile
