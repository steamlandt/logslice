// Package output provides a configurable writer for emitting sliced log lines
// to either stdout or a file. It supports multiple output formats, including
// raw line passthrough and line-number-prefixed output, making it easy to
// correlate extracted segments back to their original positions in the source
// log file.
//
// # Formats
//
// The following output formats are supported:
//
//   - [FormatRaw]: writes each line as-is, with no additional metadata.
//   - [FormatNumbered]: prefixes each line with its original line number,
//     useful for tracing output back to the source file.
//
// # Basic usage
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
//
// If Destination is empty, output is written to stdout.
package output
