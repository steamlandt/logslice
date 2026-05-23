// Package filter provides log line filtering based on level and keyword criteria.
package filter

import (
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	LevelAll   Level = iota
	LevelDebug Level = iota
	LevelInfo  Level = iota
	LevelWarn  Level = iota
	LevelError Level = iota
)

var levelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"warning": LevelWarn,
	"error": LevelError,
	"err":   LevelError,
}

// Filter holds criteria for matching log lines.
type Filter struct {
	MinLevel Level
	Keyword  string
}

// New returns a Filter with the given minimum level string and keyword.
// If levelStr is empty or unrecognised, LevelAll is used.
func New(levelStr, keyword string) *Filter {
	lvl := LevelAll
	if l, ok := levelNames[strings.ToLower(levelStr)]; ok {
		lvl = l
	}
	return &Filter{
		MinLevel: lvl,
		Keyword:  keyword,
	}
}

// Match reports whether line satisfies the filter criteria.
func (f *Filter) Match(line string) bool {
	if f.MinLevel != LevelAll {
		if !f.lineMetLevel(line) {
			return false
		}
	}
	if f.Keyword != "" {
		if !strings.Contains(line, f.Keyword) {
			return false
		}
	}
	return true
}

// lineMetLevel checks whether the line contains a level token >= f.MinLevel.
func (f *Filter) lineMetLevel(line string) bool {
	lower := strings.ToLower(line)
	for name, lvl := range levelNames {
		if lvl >= f.MinLevel && strings.Contains(lower, name) {
			return true
		}
	}
	return false
}
