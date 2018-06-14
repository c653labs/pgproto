package pgproto

import (
	"bytes"
	"fmt"
	"io"
)

// 'X' [int32 - length]
var rawTerminationMessage = [5]byte{
	// Tag
	'X',
	// Length
	'\x00', '\x00', '\x00', '\x04',
}

type Termination struct{}

func (t *Termination) client() {}

func ParseTermination(r io.Reader) (*Termination, error) {
	b := newReadBuffer(r)

	var msg [5]byte
	_, err := b.Read(msg[:])
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(msg[:], rawTerminationMessage[:]) {
		return nil, fmt.Errorf("invalid termination message")
	}
	return &Termination{}, nil
}

func (t *Termination) Encode() []byte {
	// 'X' [int32 - length]
	return rawTerminationMessage[:]
}

func (t *Termination) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "Termination",
		"Payload": nil,
	}
}

func (t *Termination) String() string { return messageToString(t) }
