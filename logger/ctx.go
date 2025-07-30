package logger

import "context"

type CtxKey struct{}

func FromCtx(ctx context.Context) (SelfLogger, bool) {
	logger, ok := ctx.Value(CtxKey{}).(SelfLogger)
	return logger, ok
}

func FromCtxOrNop(ctx context.Context) SelfLogger {
	if logger, ok := FromCtx(ctx); ok {
		return logger
	}

	return NewNop()
}

func ToCtx(ctx context.Context, lgr SelfLogger) context.Context {
	return context.WithValue(ctx, CtxKey{}, lgr)
}
