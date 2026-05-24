// Package checkpoint provides functionality for saving and resuming
// log slicing progress using persistent offset files.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// ErrNotFound is returned when no checkpoint file exists at the given path.
var ErrNotFound = errors.New("checkpoint: file not found")

// State holds the persisted progress of a slicing operation.
type State struct {
	// InputFile is the absolute path of the log file being processed.
	InputFile string `json:"input_file"`
	// ByteOffset is the byte position in the file where processing paused.
	ByteOffset int64 `json:"byte_offset"`
	// LastTimestamp is the timestamp of the last successfully processed line.
	LastTimestamp time.Time `json:"last_timestamp"`
	// SavedAt records when this checkpoint was written.
	SavedAt time.Time `json:"saved_at"`
}

// Save writes the given State to path as JSON, creating or truncating the file.
func Save(path string, s State) error {
	s.SavedAt = time.Now().UTC()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// Load reads a State from path. Returns ErrNotFound if the file does not exist.
func Load(path string) (State, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return State{}, ErrNotFound
		}
		return State{}, err
	}
	defer f.Close()
	var s State
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return State{}, err
	}
	return s, nil
}

// Remove deletes the checkpoint file at path.
// It is a no-op if the file does not exist.
func Remove(path string) error {
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
