// Package slicer extracts a time-bounded segment from a log file.
package slicer

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/example/logslice/internal/lineparser"
	"github.com/example/logslice/internal/searcher"
)

// Options controls the slice operation.
type Options struct {
	From    time.Time
	To      time.Time
	Pattern string // optional custom regex for lineparser
}

// Slice opens the log file at path and writes lines whose timestamps fall
// within [opts.From, opts.To] to dst.
func Slice(path string, opts Options, dst io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("slicer: open %s: %w", path, err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("slicer: stat %s: %w", path, err)
	}

	pattern := opts.Pattern
	if pattern == "" {
		pattern = lineparser.DefaultPattern
	}
	lp, err := lineparser.New(pattern)
	if err != nil {
		return fmt.Errorf("slicer: lineparser: %w", err)
	}

	parseFn := func(line []byte) (time.Time, error) {
		return lp.ParseTimestamp(line)
	}

	s := searcher.New(f, fi.Size(), parseFn)

	startPos, err := s.FindStart(opts.From)
	if err != nil {
		return fmt.Errorf("slicer: find start: %w", err)
	}
	endPos, err := s.FindEnd(opts.To)
	if err != nil {
		return fmt.Errorf("slicer: find end: %w", err)
	}

	if endPos <= startPos {
		return nil // empty range
	}

	if _, err := f.Seek(int64(startPos), io.SeekStart); err != nil {
		return fmt.Errorf("slicer: seek: %w", err)
	}

	length := int64(endPos) - int64(startPos)
	_, err = io.Copy(dst, io.LimitReader(f, length))
	if err != nil {
		return fmt.Errorf("slicer: copy: %w", err)
	}
	return nil
}
