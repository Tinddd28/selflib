package bufferadapter_test

import (
	"testing"

	"github.com/Tinddd28/selflib/ewrap"
	"github.com/Tinddd28/selflib/logger/bufferadapter"
	"github.com/Tinddd28/selflib/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdapter(t *testing.T) {
	t.Parallel()

	t.Run(".Log()", func(t *testing.T) {
		t.Parallel()

		buff := bufferadapter.LogEntries{}

		adapter := bufferadapter.New(&buff)

		adapter.Log(42, "foo", nil)
		require.Len(t, buff, 1)

		assert.Equal(t, bufferadapter.LogEntry{
			Level: 42,
			Msg:   "foo",
		}, buff[0])

		err := ewrap.New("some error")
		adapter.Log(42, "foo", err, types.F("foo", "bar"))
		require.Len(t, buff, 2)
		assert.Equal(t, bufferadapter.LogEntry{
			Level:  42,
			Msg:    "foo",
			Error:  err,
			Fields: types.List{types.F("foo", "bar")},
		}, buff[1])
	})

	t.Run(".WithFields()", func(t *testing.T) {
		t.Parallel()

		buff := bufferadapter.LogEntries{}

		adapterMain := bufferadapter.New(&buff)

		adapter := adapterMain.WithName("some logger").WithFields(types.F("foo", "bar"))
		require.NotSame(t, adapterMain, adapter)

		adapter.Log(42, "foo", nil)
		require.Len(t, buff, 1)
		assert.Equal(t, bufferadapter.LogEntry{
			LoggerName: "some logger",
			Level:      42,
			Msg:        "foo",
			Fields:     types.List{types.F("foo", "bar")},
		}, buff[0])

		adapter.Log(42, "foo", nil, types.F("baz", "qux"))
		require.Len(t, buff, 2)
		assert.Equal(t, bufferadapter.LogEntry{
			LoggerName: "some logger",
			Level:      42,
			Msg:        "foo",
			Fields:     types.List{types.F("foo", "bar"), types.F("baz", "qux")},
		}, buff[1])
	})

	t.Run(".WithName()", func(t *testing.T) {
		t.Parallel()

		buff := bufferadapter.LogEntries{}

		adapterMain := bufferadapter.New(&buff)

		adapter := adapterMain.WithFields(types.F("foo", "bar")).WithName("some logger")
		require.NotSame(t, adapterMain, adapter)

		adapter.Log(42, "foo", nil)
		require.Len(t, buff, 1)

		assert.Equal(t, bufferadapter.LogEntry{
			LoggerName: "some logger",
			Level:      42,
			Msg:        "foo",
			Fields:     types.List{types.F("foo", "bar")},
		}, buff[0])
	})
}
