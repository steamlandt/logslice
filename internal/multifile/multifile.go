// Package multifile provides utilities for processing multiple log files
// in sequence, merging their output into a single time-bounded stream.
package multifile

import (
	"fmt"
	"io"
	"time"

	"github.com/logslice/logslice/internal/slicer"
	"github.com/logslice/logslice/internal/config"
)

// Result holds the outcome of processing a single file.
type Result struct {
	Path    string
	Lines   int64
	Err     error
}

// Processor runs the slicer over a list of files and writes all matching
// lines to a single destination writer.
type Processor struct {
	cfg   *config.Config
	from  time.Time
	to    time.Time
}

// New creates a Processor for the given config, from, and to times.
func New(cfg *config.Config, from, to time.Time) *Processor {
	return &Processor{cfg: cfg, from: from, to: to}
}

// Run iterates over paths, slicing each file and writing matching lines to w.
// It returns one Result per file. Processing continues even if a single file
// fails so callers receive a complete picture of what succeeded.
func (p *Processor) Run(paths []string, w io.Writer) []Result {
	results := make([]Result, 0, len(paths))
	for _, path := range paths {
		n, err := p.processOne(path, w)
		results = append(results, Result{
			Path:  path,
			Lines: n,
			Err:   err,
		})
	}
	return results
}

func (p *Processor) processOne(path string, w io.Writer) (int64, error) {
	copy := *p.cfg
	copy.Input = path
	n, err := slicer.Slice(&copy, p.from, p.to, w)
	if err != nil {
		return 0, fmt.Errorf("multifile: %s: %w", path, err)
	}
	return n, nil
}
