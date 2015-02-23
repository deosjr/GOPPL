//TODO: Start by copying from csvreader, then make adjustments
// http://golang.org/src/encoding/csv/reader.go

package memory

import (
	"bufio"
	"bytes"
	"io"
	
	"GOPPL/prolog"
)

type Reader struct {
	line int
	column int
	r *bufio.Reader
	field bytes.Buffer
}

// expects UTF-8 input
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: bufio.NewReader(r),
	}
}

func (r Reader) Read() prolog.Atom {
	return prolog.Atom{"HAI"}
}