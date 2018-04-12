package pgproto

import (
	"fmt"
	"io"
)

type NoData struct{}

func (n *NoData) server() {}

func ParseNoData(r io.Reader) (*NoData, error) {
	buf := newReadBuffer(r)

	// 'n' [int32 - length]
	err := buf.ReadTag('n')
	if err != nil {
		return nil, err
	}

	_, err = buf.ReadLength()
	if err != nil {
		return nil, err
	}

	return &NoData{}, nil
}

func (n *NoData) Encode() []byte {
	// 'n' [int32 - length]
	buf := newWriteBuffer()
	buf.Wrap('n')
	return buf.Bytes()
}

func (n *NoData) WriteTo(w io.Writer) (int64, error) { return writeTo(n, w) }

func (n *NoData) String() string {
	return fmt.Sprintf("NoData<>")
}
