package searcher

import (
	"io"
	"time"
)

// Position represents a byte offset in a file.
type Position int64

// Searcher performs binary search over a sorted log file to find
// the byte offsets bounding a requested time range.
type Searcher struct {
	rs      io.ReadSeeker
	size    int64
	parseFn func(line []byte) (time.Time, error)
}

// New creates a Searcher for the given ReadSeeker whose total size is known.
// parseFn must extract a timestamp from a raw log line.
func New(rs io.ReadSeeker, size int64, parseFn func([]byte) (time.Time, error)) *Searcher {
	return &Searcher{rs: rs, size: size, parseFn: parseFn}
}

// FindStart returns the byte offset of the first line whose timestamp is >= t.
func (s *Searcher) FindStart(t time.Time) (Position, error) {
	return s.binarySearch(t, false)
}

// FindEnd returns the byte offset just past the last line whose timestamp is <= t.
func (s *Searcher) FindEnd(t time.Time) (Position, error) {
	return s.binarySearch(t, true)
}

// binarySearch locates a boundary position.
// If findEnd is false it finds the first line >= t.
// If findEnd is true it finds the position after the last line <= t.
func (s *Searcher) binarySearch(t time.Time, findEnd bool) (Position, error) {
	lo, hi := int64(0), s.size
	result := hi

	for lo < hi {
		mid := (lo + hi) / 2
		lineStart, line, err := s.readLineAt(mid)
		if err != nil {
			return 0, err
		}
		ts, err := s.parseFn(line)
		if err != nil {
			// unparseable line — advance
			lo = lineStart + int64(len(line)) + 1
			continue
		}
		if findEnd {
			if ts.After(t) {
				hi = lineStart
			} else {
				result = lineStart + int64(len(line)) + 1
				lo = result
			}
		} else {
			if ts.Before(t) {
				lo = lineStart + int64(len(line)) + 1
			} else {
				result = lineStart
				hi = lineStart
			}
		}
	}
	return Position(result), nil
}

// readLineAt seeks to approximately offset and returns the start of the
// next complete line along with its content (without the newline).
func (s *Searcher) readLineAt(offset int64) (int64, []byte, error) {
	if _, err := s.rs.Seek(offset, io.SeekStart); err != nil {
		return 0, nil, err
	}
	buf := make([]byte, 4096)
	n, err := s.rs.Read(buf)
	if err != nil && err != io.EOF {
		return 0, nil, err
	}
	buf = buf[:n]

	start := 0
	if offset > 0 {
		// skip partial line
		for start < len(buf) && buf[start] != '\n' {
			start++
		}
		start++ // skip the newline itself
	}

	end := start
	for end < len(buf) && buf[end] != '\n' {
		end++
	}
	return offset + int64(start), buf[start:end], nil
}
