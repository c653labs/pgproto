package pgproto

import (
	"fmt"
	"io"
)

type SimpleQuery struct {
	Query []byte
}

func ParseSimpleQuery(r io.Reader) (*SimpleQuery, error) {
	b := newReadBuffer(r)

	// 'Q' [int32 - length] [query] \0
	err := b.ReadTag('Q')
	if err != nil {
		return nil, err
	}

	b, err = b.ReadLength()
	if err != nil {
		return nil, err
	}

	q := &SimpleQuery{}
	q.Query, err = b.ReadString(true)
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (q *SimpleQuery) Encode() []byte {
	b := newWriteBuffer()
	b.WriteString(q.Query, true)
	b.Wrap('Q')
	return b.Bytes()
}

func (q *SimpleQuery) String() string {
	return fmt.Sprintf("SimpleQuery<Query=%#v>", string(q.Query))
}

func (q *SimpleQuery) WriteTo(w io.Writer) (int64, error) { return writeTo(q, w) }
