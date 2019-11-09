package bom

import (
	"bufio"
	"io"
)

const (
	bom0 = 0xef
	bom1 = 0xbb
	bom2 = 0xbf
)

// NewReader returns an io.Reader that will skip over initial UTF-8 byte order marks.
func NewReader(r io.Reader) io.Reader {
	buf := bufio.NewReader(r)
	b, err := buf.Peek(3)
	if err != nil {
		return buf
	}

	if b[0] == bom0 && b[1] == bom1 && b[2] == bom2 {
		buf.Discard(3)
	}

	return buf
}
