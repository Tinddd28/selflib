package types

import (
	"iter"
	"strings"
)

type List []Field

func (l List) All() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for i := 0; i < len(l); i++ {
			if !yield(l[i].K, l[i].V) {
				return
			}
		}
	}
}

func (l List) WriteTo(b *strings.Builder) {
	WriteTo(b, l.All())
}
