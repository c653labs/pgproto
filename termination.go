package pgmsg

import (
	"fmt"
	"io"
)

type Termination struct{}

func ParseTermination(r io.Reader) (*Termination, error) {
	b := NewReadBuffer(r)

	// 'X' [int32 - length]
	tag, err := b.ReadByte()
	if tag != 'X' {
		return nil, fmt.Errorf("invalid tag '%c' for termination message, msut be 'X'", tag)
	}

	l, err := b.ReadInt()
	if err != nil {
		return nil, err
	}
	if l != 4 {
		return nil, fmt.Errorf("invalid length for termination message, must be 4")
	}
	return &Termination{}, nil
}

func (t *Termination) Encode() []byte {
	// 'X' [int32 - length]
	w := NewWriteBuffer()
	w.WriteByte('X')
	w.WriteInt(4)
	return w.Bytes()
}
