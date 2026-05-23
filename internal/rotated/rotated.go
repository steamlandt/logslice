// Package rotated provides support for reading across rotated log files.
// It discovers and orders rotated log files (e.g. app.log, app.log.1, app.log.2.gz)
// so the slicer can treat them as a single continuous stream.
package rotated

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// File represents a single (possibly compressed) log file in a rotation set.
type File struct {
	Path       string
	Compressed bool
}

// Discover returns all rotated siblings of basePath in ascending age order
// (oldest first). basePath itself is included as the newest file.
//
// Example: given /var/log/app.log it may return
//
//	[app.log.3.gz, app.log.2.gz, app.log.1, app.log]
func Discover(basePath string) ([]File, error) {
	dir := filepath.Dir(basePath)
	base := filepath.Base(basePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("rotated: read dir %s: %w", dir, err)
	}

	var rotated []File
	var primary File

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == base {
			primary = File{Path: filepath.Join(dir, name), Compressed: false}
			continue
		}
		// Match siblings like app.log.1, app.log.2.gz
		stripped := strings.TrimSuffix(name, ".gz")
		if strings.HasPrefix(stripped, base+".") {
			rotated = append(rotated, File{
				Path:       filepath.Join(dir, name),
				Compressed: strings.HasSuffix(name, ".gz"),
			})
		}
	}

	// Sort rotated files: higher numeric suffix = older, so sort descending by
	// name to get oldest-first order (e.g. .3 before .2 before .1).
	sort.Slice(rotated, func(i, j int) bool {
		return rotated[i].Path > rotated[j].Path
	})

	if primary.Path != "" {
		rotated = append(rotated, primary)
	}
	return rotated, nil
}

// Open returns a ReadCloser for the given File, transparently decompressing
// gzip files.
func Open(f File) (io.ReadCloser, error) {
	fh, err := os.Open(f.Path)
	if err != nil {
		return nil, fmt.Errorf("rotated: open %s: %w", f.Path, err)
	}
	if !f.Compressed {
		return fh, nil
	}
	gr, err := gzip.NewReader(fh)
	if err != nil {
		fh.Close()
		return nil, fmt.Errorf("rotated: gzip reader %s: %w", f.Path, err)
	}
	return &gzipReadCloser{gz: gr, fh: fh}, nil
}

type gzipReadCloser struct {
	gz *gzip.Reader
	fh *os.File
}

func (g *gzipReadCloser) Read(p []byte) (int, error) { return g.gz.Read(p) }
func (g *gzipReadCloser) Close() error {
	err := g.gz.Close()
	if err2 := g.fh.Close(); err == nil {
		err = err2
	}
	return err
}
