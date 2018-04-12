package pgproto

import (
	"fmt"
	"io"
)

type Execute struct {
	Portal  []byte
	MaxRows int
}

func (e *Execute) client() {}

func ParseExecute(r io.Reader) (*Execute, error) {
	b := newReadBuffer(r)

	// 'E' [int32 - length] [string - portal] \0 [int32 - max rows]
	err := b.ReadTag('E')
	if err != nil {
		return nil, err
	}

	buf, err := b.ReadLength()
	if err != nil {
		return nil, err
	}

	e := &Execute{}

	e.Portal, err = buf.ReadString(true)
	if err != nil {
		return nil, err
	}

	e.MaxRows, err = buf.ReadInt()
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Execute) Encode() []byte {
	// 'E' [int32 - length] [string - portal] \0 [int32 - max rows]
	w := newWriteBuffer()
	w.WriteString(e.Portal, true)
	w.WriteInt(e.MaxRows)
	w.Wrap('E')
	return w.Bytes()
}

func (e *Execute) WriteTo(w io.Writer) (int64, error) { return writeTo(e, w) }

func (e *Execute) String() string {
	return fmt.Sprintf("Execute<Portal=%#v, MaxRows=%d>", string(e.Portal), e.MaxRows)
}
