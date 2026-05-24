// Package follow provides tail -f style log following for logslice.
//
// It opens a file at an optional byte offset and streams new lines to a
// channel as they are appended. Polling is used instead of inotify so the
// implementation works across all platforms and on files served over NFS or
// FUSE mounts.
//
// When the supplied context is cancelled the follower stops, reports the
// current byte offset on the provided offset channel, and closes the line
// channel. The caller can persist that offset via the checkpoint package so
// that the next run resumes exactly where the previous one left off.
//
// Typical usage:
//
//	offCh := make(chan int64, 1)
//	lines, err := follow.Lines(ctx, "/var/log/app.log",
//	    follow.Options{Offset: savedOffset}, offCh)
//	for line := range lines {
//	    process(line)
//	}
//	newOffset := <-offCh
//	checkpoint.Save(cpPath, newOffset)
package follow
