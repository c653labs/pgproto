package pgmsg

import (
	"fmt"
	"io"
)

type Message interface {
	Encode() []byte
}

func ParseMessage(r io.Reader) (Message, error) {
	// Create a buffer
	buf := NewReadBuffer(r)

	// Look at the first byte to determine the type of message we have
	start, err := buf.PeekByte()
	if err != nil {
		return nil, err
	}

	// Startup message:
	//   [int32 - length] [int32 - protocol] [[string]\0[string]\0] \0
	// Regular message
	//   [char - tag] [int32 - length] [payload] \0
	switch start {
	case '\x00':
		// StartupMessage
		return ParseStartupMessage(buf)
	case 'p':
		// Password message
		return ParsePasswordMessage(buf)
	case 'R':
		// Authentication request
		return nil, nil
	case 'S':
		// Parameter status
		return nil, nil
	case 'E':
		// Error message
		return nil, nil
	case 'X':
		// Termination
		return ParseTermination(buf)
	default:
		return nil, fmt.Errorf("unknown message tag '%c'", start)
	}
}
