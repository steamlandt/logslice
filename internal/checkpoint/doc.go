// Package checkpoint provides save/load/remove helpers for persisting
// byte-level progress through a log file.
//
// A checkpoint file is a small JSON document written alongside (or near)
// the output file. It records the source file path, the byte offset of
// the last committed line, and the timestamp of that line so that a
// subsequent run can seek directly to the right position instead of
// re-scanning from the beginning.
//
// Typical usage:
//
//	state, err := checkpoint.Load(".logslice.ckpt")
//	if errors.Is(err, checkpoint.ErrNotFound) {
//		// start from the beginning
//	}
//
//	// … process lines …
//
//	_ = checkpoint.Save(".logslice.ckpt", checkpoint.State{
//		InputFile:     cfg.Input,
//		ByteOffset:    currentOffset,
//		LastTimestamp: lastTS,
//	})
package checkpoint
