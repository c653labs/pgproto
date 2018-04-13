package pgproto

import (
	"io"
)

type CloseComplete struct{}

func (c *CloseComplete) server() {}

func ParseCloseComplete(r io.Reader) (*CloseComplete, error) {
	buf := newReadBuffer(r)

	// '3' [int32 - length]
	err := buf.ReadTag('3')
	if err != nil {
		return nil, err
	}

	_, err = buf.ReadLength()
	if err != nil {
		return nil, err
	}

	return &CloseComplete{}, nil
}

func (c *CloseComplete) Encode() []byte {
	// '3' [int32 - length]
	buf := newWriteBuffer()
	buf.Wrap('3')
	return buf.Bytes()
}

func (c *CloseComplete) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "CloseComplete",
		"Payload": nil,
	}
}

func (c *CloseComplete) String() string { return messageToString(c) }
