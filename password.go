package pgmsg

import (
	"io"
)

type PasswordMessage struct {
	Password []byte
}

func ParsePasswordMessage(r io.Reader) (*PasswordMessage, error) {
	b := newReadBuffer(r)

	// 'p' [int32 - length] [string] \0
	err := b.ReadTag('p')
	if err != nil {
		return nil, err
	}

	buf, err := b.ReadLength()
	if err != nil {
		return nil, err
	}

	p := &PasswordMessage{}

	p.Password, err = buf.ReadString(true)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *PasswordMessage) Encode() []byte {
	// 'p' [int32 - length] [string] \0
	w := newWriteBuffer()
	w.WriteString(p.Password, true)
	w.Wrap('p')
	return w.Bytes()
}

func (p *PasswordMessage) WriteTo(w io.Writer) (int, error) {
	return w.Write(p.Encode())
}
