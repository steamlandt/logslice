package config

import (
	"testing"
)

func TestFromFlagsMinimal(t *testing.T) {
	args := []string{
		"-f", "server.log",
		"-from", "2024-01-01T10:00:00Z",
		"-to", "2024-01-01T11:00:00Z",
	}
	cfg, err := FromFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.InputFile != "server.log" {
		t.Errorf("InputFile = %q, want server.log", cfg.InputFile)
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
}

func TestFromFlagsAllOptions(t *testing.T) {
	args := []string{
		"-f", "app.log",
		"-from", "2024-03-15T08:00:00Z",
		"-to", "2024-03-15T09:00:00Z",
		"-o", "out.log",
		"-n",
		"-keyword", "ERROR",
		"-level", "WARN",
		"-stats",
	}
	cfg, err := FromFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFile != "out.log" {
		t.Errorf("OutputFile = %q, want out.log", cfg.OutputFile)
	}
	if !cfg.NumberLines {
		t.Error("expected NumberLines to be true")
	}
	if cfg.Keyword != "ERROR" {
		t.Errorf("Keyword = %q, want ERROR", cfg.Keyword)
	}
	if cfg.MinLevel != "WARN" {
		t.Errorf("MinLevel = %q, want WARN", cfg.MinLevel)
	}
	if !cfg.ShowStats {
		t.Error("expected ShowStats to be true")
	}
}

func TestFromFlagsBadFrom(t *testing.T) {
	args := []string{
		"-f", "app.log",
		"-from", "not-a-date",
		"-to", "2024-01-01T11:00:00Z",
	}
	_, err := FromFlags(args)
	if err == nil {
		t.Fatal("expected error for bad --from value")
	}
}

func TestFromFlagsBadTo(t *testing.T) {
	args := []string{
		"-f", "app.log",
		"-from", "2024-01-01T10:00:00Z",
		"-to", "not-a-date",
	}
	_, err := FromFlags(args)
	if err == nil {
		t.Fatal("expected error for bad --to value")
	}
}

func TestFromFlagsUnknownFlag(t *testing.T) {
	args := []string{"-unknown-flag", "value"}
	_, err := FromFlags(args)
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}
