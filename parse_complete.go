package pgproto

import (
	"bytes"
	"fmt"
	"io"
)

// '1' [int32 - length]
var rawParseCompleteMessage = [5]byte{
	// Tag
	'1',
	// Length
	'\x00', '\x00', '\x00', '\x04',
}

type ParseComplete struct{}

func (p *ParseComplete) server() {}

func ParseParseComplete(r io.Reader) (*ParseComplete, error) {
	b := newReadBuffer(r)

	var msg [5]byte
	_, err := b.Read(msg[:])
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(msg[:], rawParseCompleteMessage[:]) {
		return nil, fmt.Errorf("invalid parse complete message")
	}

	return &ParseComplete{}, nil
}

func (p *ParseComplete) Encode() []byte {
	// '1' [int32 - length]
	return rawParseCompleteMessage[:]
}

// AsMap method returns a common map representation of this message:
//
//   map[string]interface{}{
//     "Type": "ParseComplete",
//     "Payload": nil,
//     },
//   }
func (p *ParseComplete) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "ParseComplete",
		"Payload": nil,
	}
}

func (p *ParseComplete) String() string { return messageToString(p) }
