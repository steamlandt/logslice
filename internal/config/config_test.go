package config

import (
	"testing"
	"time"
)

func base() Config {
	return Config{
		InputFile: "app.log",
		From:      time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		To:        time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
	}
}

func TestValidateOK(t *testing.T) {
	c := base()
	if err := c.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateMissingInput(t *testing.T) {
	c := base()
	c.InputFile = ""
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for missing input file")
	}
}

func TestValidateMissingFrom(t *testing.T) {
	c := base()
	c.From = time.Time{}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for zero From")
	}
}

func TestValidateMissingTo(t *testing.T) {
	c := base()
	c.To = time.Time{}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for zero To")
	}
}

func TestValidateToNotAfterFrom(t *testing.T) {
	c := base()
	c.To = c.From // equal, not after
	if err := c.Validate(); err == nil {
		t.Fatal("expected error when To == From")
	}

	c.To = c.From.Add(-time.Minute)
	if err := c.Validate(); err == nil {
		t.Fatal("expected error when To < From")
	}
}

func TestDuration(t *testing.T) {
	c := base()
	if d := c.Duration(); d != time.Hour {
		t.Fatalf("expected 1h, got %v", d)
	}
}
