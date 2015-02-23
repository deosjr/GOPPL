
package parser

import (
	"io"
	"GOPPL/types"
)

type Reader struct {
	r io.Reader
}

// expects UTF-8 input
func NewReader(r io.Reader) *Reader {
	return &Reader{r}
}

func (r Reader) Read() types.Atom {
	return types.Atom{"HAI"}
}