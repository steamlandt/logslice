package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/user/logslice/internal/timeparser"
)

// FromFlags builds a Config by parsing os.Args via the standard flag package.
// It returns an error if required flags are missing or values are invalid.
func FromFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var (
		input      = fs.String("f", "", "input log `file` (required)")
		fromStr    = fs.String("from", "", "start timestamp (required)")
		toStr      = fs.String("to", "", "end timestamp (required)")
		output     = fs.String("o", "", "output file (default: stdout)")
		numberLine = fs.Bool("n", false, "prefix output lines with line numbers")
		keyword    = fs.String("keyword", "", "only include lines containing this string")
		minLevel   = fs.String("level", "", "minimum log level (DEBUG|INFO|WARN|ERROR)")
		pattern    = fs.String("pattern", "", "custom regex with named group 'ts' for timestamp")
		showStats  = fs.Bool("stats", false, "print statistics to stderr after completion")
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &Config{
		InputFile:   *input,
		OutputFile:  *output,
		NumberLines: *numberLine,
		Keyword:     *keyword,
		MinLevel:    *minLevel,
		Pattern:     *pattern,
		ShowStats:   *showStats,
	}

	if *fromStr != "" {
		t, err := timeparser.Parse(*fromStr)
		if err != nil {
			return nil, fmt.Errorf("--from: %w", err)
		}
		cfg.From = t
	}

	if *toStr != "" {
		t, err := timeparser.Parse(*toStr)
		if err != nil {
			return nil, fmt.Errorf("--to: %w", err)
		}
		cfg.To = t
	}

	return cfg, nil
}

// MustFromFlags calls FromFlags and exits on error.
func MustFromFlags(args []string) *Config {
	cfg, err := FromFlags(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(2)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(2)
	}
	return cfg
}

// dummy to satisfy import when timeparser path differs in real module
var _ = time.RFC3339
