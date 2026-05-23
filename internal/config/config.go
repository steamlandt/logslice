// Package config handles parsing and validation of logslice CLI configuration.
package config

import (
	"errors"
	"fmt"
	"time"
)

// Config holds all runtime options for a logslice invocation.
type Config struct {
	// Input
	InputFile string

	// Time bounds
	From time.Time
	To   time.Time

	// Output
	OutputFile  string // empty means stdout
	NumberLines bool

	// Filtering
	Keyword  string
	MinLevel string

	// Line parsing
	Pattern string

	// Stats
	ShowStats bool
}

// Validate checks that the Config is internally consistent and ready for use.
func (c *Config) Validate() error {
	if c.InputFile == "" {
		return errors.New("input file must be specified")
	}
	if c.From.IsZero() {
		return errors.New("--from timestamp is required")
	}
	if c.To.IsZero() {
		return errors.New("--to timestamp is required")
	}
	if !c.To.After(c.From) {
		return fmt.Errorf("--to (%s) must be after --from (%s)", c.To.Format(time.RFC3339), c.From.Format(time.RFC3339))
	}
	return nil
}

// Duration returns the time window described by From and To.
func (c *Config) Duration() time.Duration {
	return c.To.Sub(c.From)
}
