package bufferadapter

import (
	"slices"

	"github.com/Tinddd28/selflib/logger"
	"github.com/Tinddd28/selflib/types"
)

type LogEntry struct {
	LoggerName string
	Level      int
	Msg        string
	Error      error
	Fields     types.List
}

type LogEntries []LogEntry

func (le *LogEntries) Reset() {
	*le = LogEntries{}
}

type Adapter struct {
	buff *LogEntries
	fs   types.List
	name string
}

func New(buff *LogEntries) *Adapter {
	return &Adapter{
		buff: buff,
	}
}

func (a *Adapter) Log(level int, msg string, err error, fs ...types.Field) {
	e := LogEntry{
		LoggerName: a.name,
		Level:      level,
		Msg:        msg,
		Error:      err,
	}

	e.Fields = append(slices.Clone(a.fs), fs...)
	*a.buff = append(*a.buff, e)
}

func (a *Adapter) WithFields(fs ...types.Field) logger.Adapter {
	return &Adapter{
		buff: a.buff,
		name: a.name,
		fs:   append(slices.Clone(a.fs), fs...),
	}
}

func (a *Adapter) WithName(name string) logger.Adapter {
	return &Adapter{
		buff: a.buff,
		name: name,
		fs:   slices.Clone(a.fs),
	}
}

func (a *Adapter) WithStackTrace(name string) logger.Adapter {
	return &Adapter{}
}

func (a *Adapter) Flush() error {
	return nil
}
