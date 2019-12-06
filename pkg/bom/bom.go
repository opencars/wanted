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
func NewReader(r io.Reader) (io.Reader, error) {
	buf := bufio.NewReader(r)
	b, err := buf.Peek(3)
	if err != nil {
		return nil, err
	}

	if b[0] == bom0 && b[1] == bom1 && b[2] == bom2 {
		if _, err := buf.Discard(3); err != nil {
			return nil, err
		}
	}

	return buf, nil
}
