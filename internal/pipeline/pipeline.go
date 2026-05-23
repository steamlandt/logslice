// Package pipeline wires together the core logslice components into a
// single reusable processing pipeline.
package pipeline

import (
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/config"
	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/lineparser"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/slicer"
	"github.com/yourorg/logslice/internal/stats"
)

// Pipeline holds all components needed for a single slice run.
type Pipeline struct {
	cfg    *config.Config
	parser *lineparser.Parser
	filter *filter.Filter
	writer *output.Writer
	stats  *stats.Stats
}

// New constructs a Pipeline from the supplied configuration.
// It returns an error if any component cannot be initialised.
func New(cfg *config.Config) (*Pipeline, error) {
	p, err := lineparser.New(cfg.Pattern)
	if err != nil {
		return nil, fmt.Errorf("pipeline: line parser: %w", err)
	}

	f, err := filter.New(cfg.MinLevel, cfg.Keyword)
	if err != nil {
		return nil, fmt.Errorf("pipeline: filter: %w", err)
	}

	w, err := output.New(cfg.Output, cfg.LineNumbers)
	if err != nil {
		return nil, fmt.Errorf("pipeline: output writer: %w", err)
	}

	s := stats.New(os.Stderr, cfg.ShowStats)

	return &Pipeline{
		cfg:    cfg,
		parser: p,
		filter: f,
		writer: w,
		stats:  s,
	}, nil
}

// Run executes the full slice → filter → write pipeline.
// It returns the final Stats snapshot after all processing is complete.
func (pl *Pipeline) Run() (stats.Snapshot, error) {
	pl.stats.Start()

	err := slicer.Slice(slicer.Options{
		FilePath:  pl.cfg.Input,
		From:      pl.cfg.From,
		To:        pl.cfg.To,
		Parser:    pl.parser,
		Filter:    pl.filter,
		Writer:    pl.writer,
		Stats:     pl.stats,
	})

	pl.stats.Stop()

	if err != nil {
		return stats.Snapshot{}, fmt.Errorf("pipeline: run: %w", err)
	}

	snap := pl.stats.Snapshot()
	pl.stats.Print()
	return snap, nil
}
