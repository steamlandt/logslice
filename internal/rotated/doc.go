// Package rotated discovers and opens rotated log files for a given base log
// path. Many logging systems rotate files by appending a numeric suffix
// (app.log → app.log.1 → app.log.2.gz …). This package surfaces those
// siblings so that logslice can slice across an entire rotation set as if it
// were one continuous file.
//
// # Discovery
//
// [Discover] scans the directory containing the base path and returns all
// matching siblings sorted oldest-first, with the base file (current log)
// appended last.
//
// # Opening
//
// [Open] returns an [io.ReadCloser] for any [File] returned by Discover.
// Files whose names end in ".gz" are transparently decompressed via the
// standard library's compress/gzip package.
package rotated
