package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func TestNewDefaultLevel(t *testing.T) {
	f := filter.New("", "")
	if f.MinLevel != filter.LevelAll {
		t.Fatalf("expected LevelAll, got %v", f.MinLevel)
	}
}

func TestNewUnknownLevel(t *testing.T) {
	f := filter.New("verbose", "")
	if f.MinLevel != filter.LevelAll {
		t.Fatalf("expected LevelAll for unknown level, got %v", f.MinLevel)
	}
}

func TestMatchNoFilter(t *testing.T) {
	f := filter.New("", "")
	lines := []string{
		"2024-01-01 INFO  server started",
		"2024-01-01 DEBUG connecting to db",
		"random line without level",
	}
	for _, l := range lines {
		if !f.Match(l) {
			t.Errorf("expected match for line %q", l)
		}
	}
}

func TestMatchKeyword(t *testing.T) {
	f := filter.New("", "timeout")
	if f.Match("2024-01-01 INFO  all good") {
		t.Error("expected no match without keyword")
	}
	if !f.Match("2024-01-01 ERROR connection timeout") {
		t.Error("expected match with keyword")
	}
}

func TestMatchMinLevel(t *testing.T) {
	f := filter.New("warn", "")
	tests := []struct {
		line  string
		want  bool
	}{
		{"2024-01-01 DEBUG low level noise", false},
		{"2024-01-01 INFO  informational", false},
		{"2024-01-01 WARN  disk almost full", true},
		{"2024-01-01 ERROR fatal problem", true},
	}
	for _, tc := range tests {
		got := f.Match(tc.line)
		if got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}

func TestMatchLevelAndKeyword(t *testing.T) {
	f := filter.New("error", "panic")
	if f.Match("2024-01-01 ERROR disk full") {
		t.Error("should not match: error level but missing keyword")
	}
	if f.Match("2024-01-01 WARN  panic in handler") {
		t.Error("should not match: has keyword but level too low")
	}
	if !f.Match("2024-01-01 ERROR panic in handler") {
		t.Error("should match: error level and keyword present")
	}
}

func TestMatchKeywordCaseInsensitive(t *testing.T) {
	f := filter.New("", "timeout")
	tests := []struct {
		line string
		want bool
	}{
		{"2024-01-01 ERROR connection TIMEOUT exceeded", true},
		{"2024-01-01 ERROR connection Timeout exceeded", true},
		{"2024-01-01 ERROR connection timeout exceeded", true},
		{"2024-01-01 ERROR connection error", false},
	}
	for _, tc := range tests {
		got := f.Match(tc.line)
		if got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}
