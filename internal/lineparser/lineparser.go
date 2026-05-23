package lineparser

import (
	"regexp"
	"time"

	"github.com/yourorg/logslice/internal/timeparser"
)

// Result holds the outcome of parsing a single log line.
type Result struct {
	Timestamp time.Time
	Raw       string
	Matched   bool
}

// Parser extracts timestamps from log lines using a configurable regex.
type Parser struct {
	re  *regexp.Regexp
	loc *time.Location
}

// New creates a Parser using the provided regular expression pattern.
// The pattern must contain a named capture group "ts" for the timestamp.
// If pattern is empty, a sensible default covering common log formats is used.
func New(pattern string, loc *time.Location) (*Parser, error) {
	if pattern == "" {
		pattern = `(?P<ts>\d{4}[-/]\d{2}[-/]\d{2}[T ]\d{2}:\d{2}:\d{2}(?:[.,]\d+)?(?:Z|[+-]\d{2}:?\d{2})?)`
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	if loc == nil {
		loc = time.UTC
	}
	return &Parser{re: re, loc: loc}, nil
}

// ParseLine attempts to extract a timestamp from a raw log line.
func (p *Parser) ParseLine(line string) Result {
	match := p.re.FindStringSubmatch(line)
	if match == nil {
		return Result{Raw: line, Matched: false}
	}

	tsIdx := p.re.SubexpIndex("ts")
	if tsIdx < 0 || tsIdx >= len(match) {
		return Result{Raw: line, Matched: false}
	}

	tsStr := match[tsIdx]
	t, _, err := timeparser.ParseWithLocation(tsStr, p.loc)
	if err != nil {
		return Result{Raw: line, Matched: false}
	}

	return Result{Timestamp: t, Raw: line, Matched: true}
}
