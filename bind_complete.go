package pgproto

import (
	"io"
)

type BindComplete struct{}

func (b *BindComplete) server() {}

func ParseBindComplete(r io.Reader) (*BindComplete, error) {
	buf := newReadBuffer(r)

	// '2' [int32 - length]
	err := buf.ReadTag('2')
	if err != nil {
		return nil, err
	}

	_, err = buf.ReadLength()
	if err != nil {
		return nil, err
	}

	return &BindComplete{}, nil
}

func (b *BindComplete) Encode() []byte {
	// '2' [int32 - length]
	buf := newWriteBuffer()
	buf.Wrap('2')
	return buf.Bytes()
}

func (b *BindComplete) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "BindComplete",
		"Payload": nil,
	}
}

func (b *BindComplete) String() string { return messageToString(b) }
