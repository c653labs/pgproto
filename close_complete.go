package pgproto

import (
	"bytes"
	"fmt"
	"io"
)

// '3' [int32 - length]
var rawCloseCompleteMessage = [5]byte{
	// Tag
	'3',
	// Length
	'\x00', '\x00', '\x00', '\x04',
}

// CloseComplete represents a server response message
type CloseComplete struct{}

func (c *CloseComplete) server() {}

// ParseCloseComplete will attempt to read a CloseComplete message from the io.Reader
func ParseCloseComplete(r io.Reader) (*CloseComplete, error) {
	b := newReadBuffer(r)

	var msg [5]byte
	_, err := b.Read(msg[:])
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(msg[:], rawCloseCompleteMessage[:]) {
		return nil, fmt.Errorf("invalid close complete message")
	}

	return &CloseComplete{}, nil
}

// Encode will return the byte representation of this message
func (c *CloseComplete) Encode() []byte {
	// '3' [int32 - length]
	return rawCloseCompleteMessage[:]
}

// AsMap method returns a common map representation of this message:
//
//   map[string]interface{}{
//     "Type": "CloseComplete",
//     "Payload": nil,
//   }
func (c *CloseComplete) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "CloseComplete",
		"Payload": nil,
	}
}

func (c *CloseComplete) String() string { return messageToString(c) }
