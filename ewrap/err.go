package ewrap

import (
	"errors"
	"strings"

	"github.com/Tinddd28/selflib/types"
)

type E struct {
	errs   []error
	fields types.List
}

func New(reason string, f ...types.Field) *E {
	return From(errors.New(reason), f...)
}

func NewFrom(reason string, wrapped error, f ...types.Field) *E {
	if wrapped == nil {
		return New(reason, f...)
	}
	return &E{
		errs:   []error{errors.New(reason), wrapped},
		fields: f,
	}
}

func From(origin error, f ...types.Field) *E {
	if origin == nil {
		origin = errors.New("error(nil)")
	}
	return &E{
		errs:   []error{origin},
		fields: f,
	}
}

func (e *E) Error() string {
	b := &strings.Builder{}
	writeTo(b, e)
	return b.String()
}

const separator = ": "

func writeTo(b *strings.Builder, err error) {
	if b.Len() > 0 {
		b.WriteString(separator)
	}
	ee, ok := err.(*E)
	if !ok {
		b.WriteString(err.Error())
		return
	}

	b.WriteString(ee.Reason())

	if ee == nil {
		return
	}

	if len(ee.fields) > 0 {
		b.WriteRune(' ')
		ee.fields.WriteTo(b)
	}

	if len(ee.errs) > 1 {
		writeTo(b, ee.errs[1])
	}
}

func (e *E) Reason() string {
	if e == nil {
		return "(*ewrap.E)(nil)"
	}

	if len(e.errs) == 0 {
		return "(*ewrap.E)(empty)"
	}
	return e.errs[0].Error()
}

func (e *E) Wrap(err error, f ...types.Field) *E {
	if err == nil {
		err = errors.New("error(nil)")
	}
	return &E{
		errs:   []error{e, err},
		fields: f,
	}
}

func (e *E) Unwrap() []error {
	return e.errs
}

func (e *E) WithFields(f ...types.Field) *E {
	return From(e, f...)
}

func (e *E) WithField(key string, value any) *E {
	return e.WithFields(types.F(key, value))
}

func (e *E) Fields() types.List {
	return e.fields
}

func (e *E) FindOrigin(origin error) *E {
	if e == nil {
		return nil
	}

	// Check if any direct error in the chain matches the origin
	for i, err := range e.errs {
		if errors.Is(err, origin) {
			// If it's the first error, return this E
			if i == 0 {
				return e
			}
			// Otherwise, try to cast to *E and return
			if wrappedE, ok := err.(*E); ok {
				return wrappedE
			}
			// If it's not an *E but matches, return this E as it contains the origin
			return e
		}

		// If the error is an E, recursively check it
		if wrappedE, ok := err.(*E); ok {
			err := wrappedE.FindOrigin(origin)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
