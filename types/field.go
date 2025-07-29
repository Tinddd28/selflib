package ewrap

import (
	"fmt"
	"strings"
)

type Field struct {
	K string
	V any
}

type List []Field

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
