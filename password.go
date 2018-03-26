package pgproto

import (
	"bytes"
	"fmt"
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

func (p *PasswordMessage) PasswordValid(user []byte, password []byte, salt []byte) bool {
	hash := hashPassword(user, password, salt)
	return bytes.Equal(p.Password, hash)
}

func (p *PasswordMessage) SetPassword(user []byte, password []byte, salt []byte) {
	p.Password = hashPassword(user, password, salt)
}

func (p *PasswordMessage) Encode() []byte {
	// 'p' [int32 - length] [string] \0
	w := newWriteBuffer()
	w.WriteString(p.Password, true)
	w.Wrap('p')
	return w.Bytes()
}

func (p *PasswordMessage) WriteTo(w io.Writer) (int64, error) { return writeTo(p, w) }

func (p *PasswordMessage) String() string {
	return fmt.Sprintf("PasswordMessage<Password=%#v>", string(p.Password))
}
