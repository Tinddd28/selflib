package logger

import "github.com/Tinddd28/selflib/types"

type Adapter interface {
	Log(level int, msg string, err error, fs ...types.Field)

	WithFields(fs ...types.Field) Adapter

	WithName(name string) Adapter

	WithStackTrace(trace string) Adapter

	Flush() error
}

type NopAdapter struct{}

func (NopAdapter) Log(int, string, error, ...types.Field) {}

func (a NopAdapter) WithFields(...types.Field) Adapter { return a }

func (a NopAdapter) WithName(string) Adapter { return a }

func (a NopAdapter) WithStackTrace(string) Adapter { return a }

func (NopAdapter) Flush() error { return nil }
