package timeparser_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/timeparser"
)

type parseCase struct {
	input    string
	wantYear int
	wantErr  bool
}

var parseCases = []parseCase{
	{"2024-03-15T10:22:33Z", 2024, false},
	{"2024-03-15T10:22:33.123456789Z", 2024, false},
	{"2024-03-15T10:22:33", 2024, false},
	{"2024-03-15 10:22:33", 2024, false},
	{"2024-03-15 10:22:33.000", 2024, false},
	{"2024/03/15 10:22:33", 2024, false},
	{"15/Mar/2024:10:22:33 +0000", 2024, false},
	{"not-a-timestamp", 0, true},
	{"", 0, true},
}

func TestParse(t *testing.T) {
	for _, tc := range parseCases {
		t.Run(tc.input, func(t *testing.T) {
			got, layout, err := timeparser.Parse(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got time %v (layout %q)", tc.input, got, layout)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tc.input, err)
			}
			if got.Year() != tc.wantYear {
				t.Errorf("year mismatch for %q: got %d, want %d", tc.input, got.Year(), tc.wantYear)
			}
			if layout == "" {
				t.Errorf("expected non-empty layout for input %q", tc.input)
			}
		})
	}
}

func TestParseWithLocation(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}

	got, _, err := timeparser.ParseWithLocation("2024-03-15 10:22:33", loc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, offset := got.Zone()
	if offset == 0 {
		t.Error("expected non-UTC offset for America/New_York")
	}
}

func TestKnownFormats(t *testing.T) {
	fmts := timeparser.KnownFormats()
	if len(fmts) == 0 {
		t.Error("expected at least one known format")
	}
}
