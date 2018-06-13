package pgproto

import (
	"io"
)

type ParseComplete struct{}

func (p *ParseComplete) server() {}

func ParseParseComplete(r io.Reader) (*ParseComplete, error) {
	buf := newReadBuffer(r)

	// '1' [int32 - length]
	err := buf.ReadTag('1')
	if err != nil {
		return nil, err
	}

	_, err = buf.ReadLength()
	if err != nil {
		return nil, err
	}

	return &ParseComplete{}, nil
}

func (p *ParseComplete) Encode() []byte {
	// '1' [int32 - length]
	return []byte{
		// Tag
		'1',
		// Length
		'\x00', '\x00', '\x00', '\x04',
	}
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
