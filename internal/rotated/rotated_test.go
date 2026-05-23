package rotated_test

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/rotated"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile %s: %v", name, err)
	}
}

func writeGzip(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create %s: %v", name, err)
	}
	gw := gzip.NewWriter(f)
	gw.Write([]byte(content))
	gw.Close()
	f.Close()
}

func TestDiscoverOrdering(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "app.log", "current")
	writeFile(t, dir, "app.log.1", "one")
	writeGzip(t, dir, "app.log.2.gz", "two")
	writeGzip(t, dir, "app.log.3.gz", "three")

	files, err := rotated.Discover(filepath.Join(dir, "app.log"))
	if err != nil {
		t.Fatalf("Discover: %v", err)
	}
	if len(files) != 4 {
		t.Fatalf("expected 4 files, got %d", len(files))
	}
	// Oldest first: .3.gz, .2.gz, .1, (current)
	if !endsWith(files[0].Path, "app.log.3.gz") {
		t.Errorf("files[0] should be app.log.3.gz, got %s", files[0].Path)
	}
	if !endsWith(files[3].Path, "app.log") {
		t.Errorf("files[3] should be app.log, got %s", files[3].Path)
	}
	if !files[0].Compressed {
		t.Error("files[0] should be marked compressed")
	}
	if files[3].Compressed {
		t.Error("files[3] should not be compressed")
	}
}

func TestDiscoverNonExistentDir(t *testing.T) {
	_, err := rotated.Discover("/nonexistent/path/app.log")
	if err == nil {
		t.Fatal("expected error for non-existent directory")
	}
}

func TestOpenPlain(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "app.log", "hello plain")

	f := rotated.File{Path: filepath.Join(dir, "app.log"), Compressed: false}
	rc, err := rotated.Open(f)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer rc.Close()
	data, _ := io.ReadAll(rc)
	if string(data) != "hello plain" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

func TestOpenGzip(t *testing.T) {
	dir := t.TempDir()
	writeGzip(t, dir, "app.log.1.gz", "hello gzip")

	f := rotated.File{Path: filepath.Join(dir, "app.log.1.gz"), Compressed: true}
	rc, err := rotated.Open(f)
	if err != nil {
		t.Fatalf("Open gzip: %v", err)
	}
	defer rc.Close()
	data, _ := io.ReadAll(rc)
	if string(data) != "hello gzip" {
		t.Errorf("unexpected content: %q", string(data))
	}
}

func endsWith(path, suffix string) bool {
	return filepath.Base(path) == suffix
}
