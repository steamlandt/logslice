// Package follow implements tail -f style log following with checkpoint support.
package follow

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

const (
	// DefaultPollInterval is how often the file is polled for new content.
	DefaultPollInterval = 200 * time.Millisecond
)

// Options configures the follower.
type Options struct {
	// PollInterval controls how often the file is checked for new lines.
	PollInterval time.Duration
	// Offset is the byte offset to start reading from (0 = beginning).
	Offset int64
}

// Lines follows a log file, sending new lines to the returned channel.
// The channel is closed when ctx is cancelled or an unrecoverable error occurs.
// The final byte offset is sent on offsetCh before the channel closes.
func Lines(ctx context.Context, path string, opts Options, offsetCh chan<- int64) (<-chan string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if opts.PollInterval == 0 {
		opts.PollInterval = DefaultPollInterval
	}

	if opts.Offset > 0 {
		if _, err := f.Seek(opts.Offset, io.SeekStart); err != nil {
			f.Close()
			return nil, err
		}
	}

	ch := make(chan string, 64)
	go func() {
		defer f.Close()
		defer close(ch)

		r := bufio.NewReader(f)
		ticker := time.NewTicker(opts.PollInterval)
		defer ticker.Stop()

		for {
			for {
				line, err := r.ReadString('\n')
				if len(line) > 0 {
					select {
					case ch <- line:
					case <-ctx.Done():
						sendOffset(f, offsetCh)
						return
					}
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					sendOffset(f, offsetCh)
					return
				}
			}

			select {
			case <-ctx.Done():
				sendOffset(f, offsetCh)
				return
			case <-ticker.C:
			}
		}
	}()

	return ch, nil
}

func sendOffset(f *os.File, ch chan<- int64) {
	if ch == nil {
		return
	}
	off, err := f.Seek(0, io.SeekCurrent)
	if err == nil {
		ch <- off
	}
}
