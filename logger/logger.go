package logger

import (
	"math"

	"github.com/Tinddd28/selflib/types"
)

const (
	LevelError = iota*10 + 10
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace

	DefaultLogLevel = math.MaxInt
)

type SelfLogger struct {
	maxLevel int
	adapter  Adapter
}

func New(adapter Adapter, maxLevel int) SelfLogger {
	return SelfLogger{
		maxLevel: maxLevel,
		adapter:  adapter,
	}
}

func NewNop() SelfLogger {
	return SelfLogger{
		adapter:  NopAdapter{},
		maxLevel: DefaultLogLevel,
	}
}

func (l SelfLogger) IsNop() bool {
	_, ok := l.adapter.(NopAdapter)
	return ok
}

// Error logs a message with the [LevelError] level
func (l SelfLogger) Error(msg string, err error, fs ...types.Field) {
	l.Log(LevelError, msg, err, fs...)
}

// Warning logs a message with the [LevelWarn] level
func (l SelfLogger) Warn(msg string, fs ...types.Field) {
	l.Log(LevelWarn, msg, nil, fs...)
}

// Warning logs a message with the [LevelWarn] level and the provided error
func (l SelfLogger) WarnE(msg string, err error, fs ...types.Field) {
	l.Log(LevelWarn, msg, err, fs...)
}
func (l SelfLogger) Info(msg string, fs ...types.Field) {
	l.Log(LevelInfo, msg, nil, fs...)
}
func (l SelfLogger) InfoE(msg string, err error, fs ...types.Field) {
	l.Log(LevelInfo, msg, err, fs...)
}
func (l SelfLogger) Debug(msg string, fs ...types.Field) {
	l.Log(LevelDebug, msg, nil, fs...)
}
func (l SelfLogger) DebugE(msg string, err error, fs ...types.Field) {
	l.Log(LevelDebug, msg, err, fs...)
}
func (l SelfLogger) Trace(msg string, fs ...types.Field) {
	l.Log(LevelTrace, msg, nil, fs...)
}
func (l SelfLogger) TraceE(msg string, err error, fs ...types.Field) {
	l.Log(LevelTrace, msg, err, fs...)
}

func (l SelfLogger) Log(level int, msg string, err error, fs ...types.Field) {
	if level > l.maxLevel {
		return
	}

	l.adapter.Log(level, msg, err, fs...)
}

func (l SelfLogger) WithFields(fs ...types.Field) SelfLogger {
	l.adapter = l.adapter.WithFields(fs...)
	return l
}

func (l SelfLogger) WithStackTrace(_ uint) SelfLogger {
	l.adapter = l.adapter.WithStackTrace("")
	return l
}

func (l SelfLogger) Flush() error {
	return l.adapter.Flush() //nolint:wrapcheck
}

func (l SelfLogger) WithName(name string) SelfLogger {
	l.adapter = l.adapter.WithName(name)
	return l
}
