package follow

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func writeTempLog(t *testing.T, lines []string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "follow-*.log")
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f
}

func TestLinesBasic(t *testing.T) {
	f := writeTempLog(t, []string{"line1", "line2", "line3"})
	f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	offCh := make(chan int64, 1)
	ch, err := Lines(ctx, f.Name(), Options{PollInterval: 50 * time.Millisecond}, offCh)
	if err != nil {
		t.Fatal(err)
	}

	var got []string
	for line := range ch {
		got = append(got, strings.TrimRight(line, "\n"))
		if len(got) == 3 {
			cancel()
		}
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if got[0] != "line1" || got[2] != "line3" {
		t.Fatalf("unexpected lines: %v", got)
	}
}

func TestLinesOffset(t *testing.T) {
	f := writeTempLog(t, []string{"skip", "keep"})
	info, _ := f.Stat()
	// offset past first line
	offset := int64(len("skip\n"))
	f.Close()
	_ = info

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	offCh := make(chan int64, 1)
	ch, err := Lines(ctx, f.Name(), Options{PollInterval: 50 * time.Millisecond, Offset: offset}, offCh)
	if err != nil {
		t.Fatal(err)
	}

	var got []string
	for line := range ch {
		got = append(got, strings.TrimRight(line, "\n"))
		if len(got) == 1 {
			cancel()
		}
	}

	if len(got) != 1 || got[0] != "keep" {
		t.Fatalf("expected [keep], got %v", got)
	}
}

func TestLinesFileNotFound(t *testing.T) {
	ctx := context.Background()
	_, err := Lines(ctx, "/nonexistent/path.log", Options{}, nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLinesOffsetReported(t *testing.T) {
	f := writeTempLog(t, []string{"hello"})
	f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	offCh := make(chan int64, 1)
	ch, err := Lines(ctx, f.Name(), Options{PollInterval: 50 * time.Millisecond}, offCh)
	if err != nil {
		t.Fatal(err)
	}

	<-ch
	cancel()
	// drain
	for range ch {
	}

	select {
	case off := <-offCh:
		if off <= 0 {
			t.Fatalf("expected positive offset, got %d", off)
		}
	case <-time.After(time.Second):
		t.Fatal("offset not reported")
	}
}
