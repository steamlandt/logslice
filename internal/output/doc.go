// Package output provides a configurable writer for emitting sliced log lines
// to either stdout or a file. It supports multiple output formats, including
// raw line passthrough and line-number-prefixed output, making it easy to
// correlate extracted segments back to their original positions in the source
// log file.
//
// Basic usage:
//
//	w, err := output.New(output.Options{
//		Format:      output.FormatNumbered,
//		Destination: "/tmp/slice.log",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer w.Close()
//
//	w.WriteLine(lineNum, lineContent)
package output
