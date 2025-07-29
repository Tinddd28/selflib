package ewrap

import (
	"errors"
	"strings"
)

type E struct {
	err     error
	wrapped error
	fields  List
}

func New(reason string, f ...Field) *E {
	return &E{
		err:     errors.New(reason),
		wrapped: nil,
		fields:  f,
	}
}

func NewFrom(origin error, reason string, f ...Field) *E {
	return &E{
		err:     errors.New(reason),
		wrapped: origin,
		fields:  f,
	}
}

func From(origin error) *E {
	return &E{
		err:     origin,
		wrapped: nil,
		fields:  nil,
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

	if ee.err == nil {
		b.WriteString("nil")
		return
	}

	b.WriteString(ee.err.Error())

	if ee.fields != nil {

	}

	if ee.wrapped != nil {
		writeTo(b, ee.wrapped)
	}
}
