package checkpoint_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "test.ckpt")
}

func TestSaveAndLoad(t *testing.T) {
	path := tempPath(t)
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	orig := checkpoint.State{
		InputFile:     "/var/log/app.log",
		ByteOffset:    4096,
		LastTimestamp: ts,
	}

	if err := checkpoint.Save(path, orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := checkpoint.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if got.InputFile != orig.InputFile {
		t.Errorf("InputFile: got %q, want %q", got.InputFile, orig.InputFile)
	}
	if got.ByteOffset != orig.ByteOffset {
		t.Errorf("ByteOffset: got %d, want %d", got.ByteOffset, orig.ByteOffset)
	}
	if !got.LastTimestamp.Equal(orig.LastTimestamp) {
		t.Errorf("LastTimestamp: got %v, want %v", got.LastTimestamp, orig.LastTimestamp)
	}
	if got.SavedAt.IsZero() {
		t.Error("SavedAt should be set by Save")
	}
}

func TestLoadNotFound(t *testing.T) {
	_, err := checkpoint.Load(filepath.Join(t.TempDir(), "missing.ckpt"))
	if !errors.Is(err, checkpoint.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestRemove(t *testing.T) {
	path := tempPath(t)
	if err := checkpoint.Save(path, checkpoint.State{InputFile: "x"}); err != nil {
		t.Fatal(err)
	}
	if err := checkpoint.Remove(path); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Error("file should have been deleted")
	}
}

func TestRemoveNonExistent(t *testing.T) {
	// Should not return an error for a missing file.
	if err := checkpoint.Remove(filepath.Join(t.TempDir(), "ghost.ckpt")); err != nil {
		t.Errorf("Remove on missing file: %v", err)
	}
}

func TestSaveInvalidPath(t *testing.T) {
	err := checkpoint.Save("/nonexistent_dir/sub/file.ckpt", checkpoint.State{})
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
