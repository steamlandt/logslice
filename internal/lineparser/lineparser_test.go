package lineparser_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/lineparser"
)

func mustParser(t *testing.T, pattern string) *lineparser.Parser {
	t.Helper()
	p, err := lineparser.New(pattern, time.UTC)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	return p
}

func TestParseLineDefaultPattern(t *testing.T) {
	p := mustParser(t, "")

	cases := []struct {
		line        string
		wantMatched bool
		wantYear    int
	}{
		{`2024-03-15T10:22:33Z INFO starting server`, true, 2024},
		{`2024-03-15 10:22:33 ERROR disk full`, true, 2024},
		{`no timestamp here`, false, 0},
		{``, false, 0},
	}

	for _, tc := range cases {
		t.Run(tc.line, func(t *testing.T) {
			res := p.ParseLine(tc.line)
			if res.Matched != tc.wantMatched {
				t.Errorf("Matched=%v, want %v for line %q", res.Matched, tc.wantMatched, tc.line)
			}
			if res.Raw != tc.line {
				t.Errorf("Raw mismatch: got %q, want %q", res.Raw, tc.line)
			}
			if tc.wantMatched && res.Timestamp.Year() != tc.wantYear {
				t.Errorf("year=%d, want %d", res.Timestamp.Year(), tc.wantYear)
			}
		})
	}
}

func TestParseLineCustomPattern(t *testing.T) {
	// Pattern for nginx-style: [15/Mar/2024:10:22:33 +0000]
	p := mustParser(t, `\[(?P<ts>\d{2}/\w+/\d{4}:\d{2}:\d{2}:\d{2} [+-]\d{4})\]`)

	line := `192.168.1.1 - - [15/Mar/2024:10:22:33 +0000] "GET / HTTP/1.1" 200 1234`
	res := p.ParseLine(line)
	if !res.Matched {
		t.Fatalf("expected match for nginx log line")
	}
	if res.Timestamp.Year() != 2024 {
		t.Errorf("year=%d, want 2024", res.Timestamp.Year())
	}
}

func TestNewInvalidPattern(t *testing.T) {
	_, err := lineparser.New(`[invalid`, time.UTC)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}
