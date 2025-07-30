package slogadapter

import (
	"context"
	"log/slog"

	"github.com/Tinddd28/selflib/logger"
	"github.com/Tinddd28/selflib/types"
)

type Adapter struct {
	lgr *slog.Logger
}

func New(lgr *slog.Logger) *Adapter {
	return &Adapter{
		lgr: lgr,
	}
}

func (a *Adapter) Log(level int, msg string, err error, fs ...types.Field) {
	sl := slog.LevelInfo

	switch level {
	case logger.LevelError:
		sl = slog.LevelError
	case logger.LevelWarn:
		sl = slog.LevelWarn
	case logger.LevelInfo:
		sl = slog.LevelInfo
	case logger.LevelDebug:
		sl = slog.LevelDebug
	case logger.LevelTrace:
		sl = slog.LevelDebug
	default:
		a.lgr.Error("unknown log level", slog.Int("got-level", level))
	}

	a.lgr.LogAttrs(context.Background(), sl, msg, fieldsListToSlogAttrs(fs, err)...)
}

func (a *Adapter) WithFields(fs ...types.Field) logger.Adapter {
	return &Adapter{
		lgr: slog.New(a.lgr.Handler().WithAttrs(fieldsListToSlogAttrs(fs, nil))),
	}
}

func (a *Adapter) WithName(name string) logger.Adapter {
	return &Adapter{
		lgr: a.lgr.WithGroup(name),
	}
}

func (a *Adapter) Flush() error {
	return nil
}

func (a *Adapter) WithStackTrace(_ string) logger.Adapter {
	return a
}

func fieldsListToSlogAttrs(fs types.List, err error) []slog.Attr {
	sfs := make([]slog.Attr, 0, len(fs)+boolToInt(err != nil))

	if err != nil {
		sfs = append(sfs, slog.Any("error", err)) // FIXME: Implement usage error; now: cant display error, only key = @e
	}

	for _, f := range fs {
		sfs = append(sfs, slog.Any(f.K, f.V))
	}

	return sfs
}

func boolToInt(v bool) int {
	if v {
		return 1
	}

	return 0
}
