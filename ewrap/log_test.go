package ewrap_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Tinddd28/selflib/ewrap"
	"github.com/Tinddd28/selflib/logger"
	"github.com/Tinddd28/selflib/logger/bufferadapter"
	"github.com/Tinddd28/selflib/types"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	t.Parallel()

	e1 := errors.New("e1")
	e2 := fmt.Errorf("e2: %w", e1)
	e3 := ewrap.NewFrom("e3", e1, types.F("foo", "bar"))
	e4 := ewrap.New("e4")

	entries := make(bufferadapter.LogEntries, 0)
	log := logger.New(bufferadapter.New(&entries), logger.DefaultLogLevel)

	ewrap.Log(nil, log.Error)
	ewrap.Log(e1, log.Error)
	ewrap.Log(e1, log.WarnE)
	ewrap.Log(e2, log.InfoE)
	ewrap.Log(e3, log.DebugE)
	ewrap.Log(e4, log.TraceE)

	assert.Equal(t, bufferadapter.LogEntries{
		{
			LoggerName: "",
			Level:      logger.LevelError,
			Msg:        e1.Error(),
			Error:      nil,
			Fields:     nil,
		},
		{
			LoggerName: "",
			Level:      logger.LevelWarn,
			Msg:        e1.Error(),
			Error:      nil,
			Fields:     nil,
		},
		{
			LoggerName: "",
			Level:      logger.LevelInfo,
			Msg:        e2.Error(),
			Error:      nil,
			Fields:     nil,
		},
		{
			LoggerName: "",
			Level:      logger.LevelDebug,
			Msg:        e3.Reason(),
			Error:      e1,
			Fields:     e3.Fields(),
		},
		{
			LoggerName: "",
			Level:      logger.LevelTrace,
			Msg:        e4.Reason(),
			Error:      nil,
			Fields:     nil,
		},
	}, entries)

}
