package pgmsg

import (
	"bytes"
	"fmt"
	"io"
)

type PasswordMessage struct {
	Password []byte
}

func ParsePasswordMessage(r io.Reader) (*PasswordMessage, error) {
	b := NewReadBuffer(r)

	// 'p' [int32 - length] [string] \0
	tag, err := b.ReadByte()
	if tag != 'p' {
		return nil, fmt.Errorf("invalid tag '%c' for password message, must be 'p'", tag)
	}

	_, raw, err := b.ReadLength()
	if err != nil {
		return nil, err
	}

	// Replace the passed in buffer with one that is only scoped to the desired length we need
	b = NewReadBuffer(bytes.NewReader(raw))

	p := &PasswordMessage{}

	p.Password, err = b.ReadString()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *PasswordMessage) Encode() []byte {
	// 'p' [int32 - length] [string] \0
	w := NewWriteBuffer()
	w.WriteString(p.Password)
	w.PrependLength()
	w.PrependByte('p')

	return w.Bytes()

}
