package slogadapter_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"slices"
	"testing"

	"github.com/Tinddd28/selflib/logger"
	"github.com/Tinddd28/selflib/logger/slogadapter"
	"github.com/Tinddd28/selflib/types"
)

type bufferEntry struct {
	Message string
	Level   slog.Level
	Attrs   map[string]any
	Group   string
}

type entriesBuffer []bufferEntry

func (b *entriesBuffer) Reset() {
	*b = []bufferEntry{}
}

type bufferHandler struct {
	buf   *entriesBuffer
	attrs []slog.Attr
	group string
}

func (*bufferHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *bufferHandler) Handle(_ context.Context, r slog.Record) error {
	entry := bufferEntry{
		Message: r.Message,
		Level:   r.Level,
		Attrs:   map[string]any{},
		Group:   h.group,
	}

	for _, a := range h.attrs {
		entry.Attrs[a.Key] = a.Value.Any()
	}

	r.Attrs(func(a slog.Attr) bool {
		entry.Attrs[a.Key] = a.Value.Any()

		return true
	})

	*h.buf = append(*h.buf, entry)

	return nil
}

func (h *bufferHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &bufferHandler{
		buf:   h.buf,
		attrs: append(slices.Clone(h.attrs), attrs...),
		group: h.group,
	}
}

func (h *bufferHandler) WithGroup(name string) slog.Handler {
	return &bufferHandler{
		buf:   h.buf,
		attrs: h.attrs,
		group: name,
	}
}

func newAdapter(buf *entriesBuffer) *slogadapter.Adapter {
	return slogadapter.New(slog.New(&bufferHandler{buf: buf})) //nolint:exhaustruct
}

const errorKey = "error"

func TestAdapter(t *testing.T) {
	t.Parallel()

	t.Run("Log", func(t *testing.T) {
		t.Parallel()

		buf := entriesBuffer{}
		adapter := newAdapter(&buf)

		adapter.Log(42, "foo", nil)

		require.Len(t, buf, 2)
		assert.Equal(t, slog.LevelError, buf[0].Level)
		assert.Nil(t, buf[0].Attrs[errorKey])
		assert.Equal(t, slog.LevelInfo, buf[1].Level)
		assert.Nil(t, buf[1].Attrs[errorKey])

		buf.Reset()

		tt := []struct {
			level     int
			slogLevel slog.Level
			msg       string
			err       error
		}{
			{logger.LevelDebug, slog.LevelDebug, "debug", errors.New("debug")},
			{logger.LevelInfo, slog.LevelInfo, "info", errors.New("info")},
			{logger.LevelWarn, slog.LevelWarn, "warning", errors.New("warning")},
			{logger.LevelError, slog.LevelError, "error", errors.New("error")},
			{logger.LevelTrace, slog.LevelDebug, "trace", errors.New("trace")},
		}

		for _, tc := range tt {
			buf.Reset()

			adapter.Log(tc.level, tc.msg, tc.err)

			require.Len(t, buf, 1)
			assert.Equal(t, tc.slogLevel, buf[0].Level)
			assert.Equal(t, tc.msg, buf[0].Message)
			assert.Same(t, tc.err, buf[0].Attrs[errorKey])
		}

	})

	t.Run("WithFields", func(t *testing.T) {
		t.Parallel()

		buf := entriesBuffer{}
		adapter := newAdapter(&buf).WithFields(types.F("foo", "bar"), types.F("baz", 42))

		adapter.Log(logger.LevelInfo, "test", nil)
		adapter.Log(logger.LevelDebug, "test", nil, types.F("bux", "qux"))
		adapter.Log(logger.LevelDebug, "test", nil, types.F("baz", "bar"))

		require.Len(t, buf, 3)

		assert.Equal(t, "bar", buf[0].Attrs["foo"])
		assert.Equal(t, int64(42), buf[0].Attrs["baz"])

		assert.Equal(t, "bar", buf[1].Attrs["foo"])
		assert.Equal(t, int64(42), buf[1].Attrs["baz"])
		assert.Equal(t, "qux", buf[1].Attrs["bux"])

		assert.Equal(t, "bar", buf[2].Attrs["foo"])
		assert.Equal(t, "bar", buf[2].Attrs["baz"])
	})

	t.Run("WithName", func(t *testing.T) {
		t.Parallel()

		buf := entriesBuffer{}
		adapter := newAdapter(&buf).WithName("test-logger")

		adapter.Log(logger.LevelInfo, "test", nil)

		require.Len(t, buf, 1)

		assert.Equal(t, "test-logger", buf[0].Group)

	})
}
