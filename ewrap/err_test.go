package ewrap_test

import (
	"errors"
	"github.com/Tinddd28/selflib/ewrap"
	"github.com/Tinddd28/selflib/types"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestErr(t *testing.T) {
	t.Parallel()

	t.Run(".Error()", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			name     string
			in       *ewrap.E
			expected string
		}{
			{
				name:     "nil",
				in:       nil,
				expected: "(*ewrap.E)(nil)",
			},
			{
				name:     "empty",
				in:       &ewrap.E{},
				expected: "(*ewrap.E)(empty)",
			},
			{
				name:     "from nil error",
				in:       ewrap.From(nil),
				expected: "error(nil)",
			},
			{
				name:     "simple",
				in:       ewrap.New("reason"),
				expected: "reason",
			},
			{
				name:     "error wrapped",
				in:       ewrap.NewFrom("error", errors.New("wrapped")),
				expected: "error: wrapped",
			},
			{
				name:     "with fields",
				in:       ewrap.New("reason", types.F("key", "value")),
				expected: "reason (key=value)",
			},
			{
				name:     "wrapped error with fields",
				in:       ewrap.NewFrom("error", errors.New("wrapped"), types.F("k1", "v1"), types.F("k2", "v2")),
				expected: "error (k1=v1, k2=v2): wrapped",
			},
			{
				name:     "nil error in the middle",
				in:       ewrap.New("e1").Wrap((*ewrap.E)(nil)).Wrap(ewrap.New("e2")),
				expected: "e1: (*ewrap.E)(nil): e2",
			},
			{
				name:     "from external error with fields",
				in:       ewrap.From(errors.New("error"), types.F("key", "value")),
				expected: "error (key=value)",
			},
			{
				name:     "some isdugsrg",
				in:       ewrap.From(ewrap.NewFrom("error", ewrap.New("wrapped", types.F("k1", "v1"))), types.F("k2", "v2")),
				expected: "error: wrapped (k1=v1) (k2=v2)",
			},
		}

		for _, tc := range tt {
			log.Printf("expected: %s;\t input: %s", tc.expected, tc.in.Error())
			assert.Equal(t, tc.expected, tc.in.Error(), tc.name)
		}

	})

	t.Run(".Wrap()", func(t *testing.T) {
		t.Parallel()

		e1 := ewrap.New("e1", types.F("k1", "v1"))
		e2 := ewrap.New("e2", types.F("k2", "v2"))

		assert.NotSame(t, e1, e1.Wrap(e2))
		assert.NotSame(t, e2, e2.Wrap(e1))
		assert.NotErrorAs(t, errors.Unwrap(e1.Wrap(e2)), "errors unwrap returns nil")

		assert.Equal(t, "e1 (k1=v1): e2 (k2=v2)", e1.Wrap(e2).Error())
		assert.Equal(t, "e1 (k1=v1) (k3=v3): e2 (k2=v2)", e1.Wrap(e2, types.F("k3", "v3")).Error())

		assert.Equal(t, "e1 (k1=v1): error(nil)", e1.Wrap(nil).Error())
	})

	t.Run(".Is()", func(t *testing.T) {
		t.Parallel()

		var (
			e0       = errors.New("e0")
			e1 error = ewrap.NewFrom("e1", os.ErrNotExist)
			e2 error = ewrap.From(e0)
			e3 error = ewrap.NewFrom("e3", e1)
			e4       = ewrap.From(e0)
		)

		assert.ErrorIs(t, e1, e1)
		assert.NotErrorIs(t, e1, e0)

		assert.ErrorIs(t, e2, e0)

		assert.ErrorIs(t, e3, e1)
		assert.ErrorIs(t, e3, os.ErrNotExist)

		assert.NotErrorIs(t, e4, e2)
		log.Println(e4.WithField("k1", "v1"))
		assert.NotErrorIs(t, e4.WithField("k1", "v1"), e4.WithField("k1", "v1"))
	})

}
