package types

import (
	"fmt"
	"iter"
	"strings"
)

type Field struct {
	K string
	V any
}

func writeKVTo(b *strings.Builder, key string, value any) {
	b.WriteString(key)
	b.WriteRune('=')

	switch val := value.(type) {
	case string:
		b.WriteString(val)
	case fmt.Stringer:
		b.WriteString(val.String())
	case error:
		b.WriteString(val.Error())
	default:
		_, _ = fmt.Fprintf(b, "%v", val)
	}
}

func (f Field) WriteTo(b *strings.Builder) {
	writeKVTo(b, f.K, f.V)
}

func (f *Field) String() string {
	b := &strings.Builder{}
	f.WriteTo(b)
	return b.String()
}

func WriteTo(b *strings.Builder, seq iter.Seq2[string, any]) {
	first := true

	for k, v := range seq {
		if first {
			first = false

			b.WriteString("(")
		} else {
			b.WriteString(", ")
		}

		writeKVTo(b, k, v)
	}

	if !first {
		b.WriteString(")")
	}
}

func F(key string, value any) Field {
	return Field{K: key, V: value}
}
