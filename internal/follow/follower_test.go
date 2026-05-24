package follow

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/checkpoint"
	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/lineparser"
)

func mustParser(t *testing.T) *lineparser.Parser {
	t.Helper()
	p, err := lineparser.New("")
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func mustFilter(t *testing.T) *filter.Filter {
	t.Helper()
	f, err := filter.New("", "")
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestFollowerRunBasic(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "app.log")
	cpPath := filepath.Join(dir, "app.cp")

	lines := []string{
		`2024-01-15 10:00:00 INFO hello world`,
		`2024-01-15 10:00:01 DEBUG debug msg`,
	}
	f, _ := os.Create(logPath)
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	fl := NewFollower(logPath, cpPath, mustParser(t), mustFilter(t), Options{
		PollInterval: 50 * time.Millisecond,
	})

	var buf bytes.Buffer
	// cancel quickly so we don't block forever
	go func() {
		time.Sleep(300 * time.Millisecond)
		cancel()
	}()

	if err := fl.Run(ctx, &buf); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "hello world") {
		t.Errorf("expected 'hello world' in output, got: %q", out)
	}
}

func TestFollowerCheckpointResume(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "app.log")
	cpPath := filepath.Join(dir, "app.cp")

	// Write two lines, save checkpoint after first.
	line1 := "2024-01-15 10:00:00 INFO first\n"
	line2 := "2024-01-15 10:00:01 INFO second\n"
	os.WriteFile(logPath, []byte(line1+line2), 0o644)
	_ = checkpoint.Save(cpPath, int64(len(line1)))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	go func() {
		time.Sleep(300 * time.Millisecond)
		cancel()
	}()

	fl := NewFollower(logPath, cpPath, mustParser(t), mustFilter(t), Options{
		PollInterval: 50 * time.Millisecond,
	})

	var buf bytes.Buffer
	fl.Run(ctx, &buf)

	if strings.Contains(buf.String(), "first") {
		t.Error("expected 'first' to be skipped due to checkpoint offset")
	}
	if !strings.Contains(buf.String(), "second") {
		t.Errorf("expected 'second' in output, got: %q", buf.String())
	}
}
