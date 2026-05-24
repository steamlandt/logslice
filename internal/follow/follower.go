package follow

import (
	"context"
	"io"

	"github.com/user/logslice/internal/checkpoint"
	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/lineparser"
)

// Follower combines file following with line parsing, filtering, and
// checkpoint persistence.
type Follower struct {
	path      string
	cpPath    string
	parser    *lineparser.Parser
	filter    *filter.Filter
	opts      Options
}

// NewFollower creates a Follower for the given log file.
// cpPath is the checkpoint file path used to persist the read offset.
func NewFollower(path, cpPath string, p *lineparser.Parser, f *filter.Filter, opts Options) *Follower {
	return &Follower{
		path:   path,
		cpPath: cpPath,
		parser: p,
		filter: f,
		opts:   opts,
	}
}

// Run starts following the log file, writing matching lines to w.
// It blocks until ctx is cancelled. The byte offset is saved to the
// checkpoint file on exit.
func (fl *Follower) Run(ctx context.Context, w io.Writer) error {
	offset, _ := checkpoint.Load(fl.cpPath)

	if fl.opts.Offset == 0 && offset > 0 {
		fl.opts.Offset = offset
	}

	offCh := make(chan int64, 1)
	ch, err := Lines(ctx, fl.path, fl.opts, offCh)
	if err != nil {
		return err
	}

	for raw := range ch {
		entry, err := fl.parser.ParseLine(raw)
		if err != nil {
			continue
		}
		if !fl.filter.Match(entry) {
			continue
		}
		if _, err := io.WriteString(w, raw); err != nil {
			break
		}
	}

	if newOffset, ok := <-offCh; ok && newOffset > 0 {
		_ = checkpoint.Save(fl.cpPath, newOffset)
	}
	return nil
}
