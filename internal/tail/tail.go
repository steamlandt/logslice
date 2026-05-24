// Package tail provides functionality to read the last N lines of a log file
// efficiently using a backward-scanning strategy without loading the full file.
package tail

import (
	"bytes"
	"errors"
	"io"
	"os"
)

const defaultChunkSize = 4096

// Lines reads the last n lines from the file at the given path.
// It scans the file backwards in chunks to avoid loading the entire file
// into memory. Returns an error if the file cannot be opened or read.
func Lines(path string, n int) ([]string, error) {
	if n <= 0 {
		return nil, errors.New("tail: n must be greater than zero")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return readLastLines(f, n)
}

// readLastLines performs the backward scan on any ReadSeeker.
func readLastLines(rs io.ReadSeeker, n int) ([]string, error) {
	size, err := rs.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	var collected [][]byte
	remaining := size
	var leftover []byte

	for len(collected) <= n && remaining > 0 {
		chunk := int64(defaultChunkSize)
		if chunk > remaining {
			chunk = remaining
		}
		remaining -= chunk

		if _, err := rs.Seek(remaining, io.SeekStart); err != nil {
			return nil, err
		}

		buf := make([]byte, chunk)
		if _, err := io.ReadFull(rs, buf); err != nil {
			return nil, err
		}

		buf = append(buf, leftover...)
		parts := bytes.Split(buf, []byte("\n"))
		leftover = parts[0]

		for i := len(parts) - 1; i >= 1; i-- {
			if len(parts[i]) == 0 {
				continue
			}
			collected = append(collected, parts[i])
			if len(collected) == n {
				break
			}
		}
	}

	if len(leftover) > 0 && len(collected) < n {
		collected = append(collected, leftover)
	}

	// Reverse so lines are in chronological order.
	for i, j := 0, len(collected)-1; i < j; i, j = i+1, j-1 {
		collected[i], collected[j] = collected[j], collected[i]
	}

	result := make([]string, len(collected))
	for i, b := range collected {
		result[i] = string(b)
	}
	return result, nil
}
