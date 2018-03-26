package pgproto

import (
	"fmt"
	"io"
)

type Termination struct{}

func ParseTermination(r io.Reader) (*Termination, error) {
	b := newReadBuffer(r)

	// 'X' [int32 - length]
	err := b.ReadTag('X')
	if err != nil {
		return nil, err
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
	w := newWriteBuffer()
	w.Wrap('X')
	return w.Bytes()
}

func (t *Termination) WriteTo(w io.Writer) (int, error) {
	return w.Write(t.Encode())
}

func (t *Termination) String() string {
	return "Termination<>"
}
