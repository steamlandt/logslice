package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Format defines the output format for sliced log lines.
type Format int

const (
	// FormatRaw writes lines as-is.
	FormatRaw Format = iota
	// FormatNumbered prepends each line with its original line number.
	FormatNumbered
)

// Options configures the Writer.
type Options struct {
	Format     Format
	Destination string // empty means stdout
}

// Writer writes log lines to a destination.
type Writer struct {
	opts   Options
	out    io.Writer
	closer io.Closer
}

// New creates a new Writer based on the provided options.
func New(opts Options) (*Writer, error) {
	var out io.Writer
	var closer io.Closer

	if opts.Destination == "" {
		out = os.Stdout
	} else {
		f, err := os.Create(opts.Destination)
		if err != nil {
			return nil, fmt.Errorf("output: create file %q: %w", opts.Destination, err)
		}
		out = f
		closer = f
	}

	return &Writer{
		opts:   opts,
		out:    bufio.NewWriter(out),
		closer: closer,
	}, nil
}

// WriteLine writes a single log line, optionally decorated with its line number.
func (w *Writer) WriteLine(lineNum int, line string) error {
	var err error
	switch w.opts.Format {
	case FormatNumbered:
		_, err = fmt.Fprintf(w.out, "%d\t%s\n", lineNum, line)
	default:
		_, err = fmt.Fprintf(w.out, "%s\n", line)
	}
	return err
}

// Flush flushes any buffered data to the underlying writer.
func (w *Writer) Flush() error {
	if bw, ok := w.out.(*bufio.Writer); ok {
		return bw.Flush()
	}
	return nil
}

// Close flushes and closes the writer if it owns the underlying file.
func (w *Writer) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}
	if w.closer != nil {
		return w.closer.Close()
	}
	return nil
}
