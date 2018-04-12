package pgproto

import (
	"fmt"
	"io"
)

type ParseComplete struct{}

func (p *ParseComplete) server() {}

func ParseParseComplete(r io.Reader) (*ParseComplete, error) {
	buf := newReadBuffer(r)

	// '1' [int32 - length]
	err := buf.ReadTag('1')
	if err != nil {
		return nil, err
	}

	_, err = buf.ReadLength()
	if err != nil {
		return nil, err
	}

	return &ParseComplete{}, nil
}

func (p *ParseComplete) Encode() []byte {
	// '1' [int32 - length]
	buf := newWriteBuffer()
	buf.Wrap('1')
	return buf.Bytes()
}

func (p *ParseComplete) WriteTo(w io.Writer) (int64, error) { return writeTo(p, w) }

func (p *ParseComplete) String() string {
	return fmt.Sprintf("ParseComplete<>")
}
