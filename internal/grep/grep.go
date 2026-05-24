// Package grep provides regex-based line scanning across log segments.
package grep

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

// Result holds a matched line along with its 1-based line number.
type Result struct {
	LineNo int
	Line   string
}

// Options controls the behaviour of Scan.
type Options struct {
	// Pattern is the regular expression to match against each line.
	Pattern string
	// MaxResults limits the number of results returned. 0 means unlimited.
	MaxResults int
	// Invert returns lines that do NOT match the pattern.
	Invert bool
}

// Scan reads all lines from r and returns those matching opts.Pattern.
// It returns an error if the pattern is invalid or an I/O error occurs.
func Scan(r io.Reader, opts Options) ([]Result, error) {
	if opts.Pattern == "" {
		return nil, fmt.Errorf("grep: pattern must not be empty")
	}
	re, err := regexp.Compile(opts.Pattern)
	if err != nil {
		return nil, fmt.Errorf("grep: invalid pattern %q: %w", opts.Pattern, err)
	}

	var results []Result
	scanner := bufio.NewScanner(r)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		text := scanner.Text()
		matched := re.MatchString(text)
		if opts.Invert {
			matched = !matched
		}
		if !matched {
			continue
		}
		results = append(results, Result{LineNo: lineNo, Line: text})
		if opts.MaxResults > 0 && len(results) >= opts.MaxResults {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("grep: scan error: %w", err)
	}
	return results, nil
}
