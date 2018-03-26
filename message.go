package pgmsg

import (
	"fmt"
	"io"
)

type Message interface {
	Encode() []byte
	WriteTo(w io.Writer) (int, error)
	String() string
}

func ParseMessage(r io.Reader) (Message, error) {
	// Create a buffer
	buf := newReadBuffer(r)

	// Look at the first byte to determine the type of message we have
	start, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}

	// Startup message:
	//   [int32 - length] [int32 - protocol] [[string]\0[string]\0] \0
	// Regular message
	//   [char - tag] [int32 - length] [payload] \0
	switch start {
	// TODO: We need to handle this case better, it might not always start with \x00
	case '\x00':
		// [int32 - length] [payload]
		// StartupMessage
		// Read the next 3 bytes, prepend with the 1 we already read to parse the length from this message
		b := make([]byte, 3)
		_, err := buf.Read(b)
		if err != nil {
			return nil, err
		}
		b = append([]byte{start}, b...)
		l := bytesToInt(b)

		// Read the rest of the message into a []byte
		// DEV: Subtract 4 to account for the length of the in32 we just read
		b = make([]byte, l-4)
		_, err = buf.Read(b)
		if err != nil {
			return nil, err
		}

		// Rebuild the message into a []byte
		w := newWriteBuffer()
		w.WriteInt(l)
		w.WriteBytes(b)

		return ParseStartupMessage(w.Reader())
	default:
		// [char tag] [int32 length] [payload]
		// Parse length from the message
		l, err := buf.ReadInt()
		if err != nil {
			return nil, err
		}

		// Read the rest of the message into a []byte
		// DEV: Subtract 4 to account for the length of the int32 we just read
		b := make([]byte, l-4)
		_, err = buf.Read(b)
		if err != nil {
			return nil, err
		}

		// Rebuild the message into a []byte
		w := newWriteBuffer()
		w.WriteByte(start)
		w.WriteInt(l)
		w.WriteBytes(b)

		switch start {
		case 'p':
			// Password message
			return ParsePasswordMessage(w.Reader())
		case 'R':
			// Authentication request
			return ParseAuthenticationRequest(w.Reader())
		case 'S':
			// Parameter status
			return ParseParameterStatus(w.Reader())
		case 'K':
			// Backend key data
			return ParseBackendKeyData(w.Reader())
		case 'Z':
			// Ready for query
			return ParseReadyForQuery(w.Reader())
		case 'Q':
			// Query
			return ParseReadyForQuery(w.Reader())
		case 'E':
			// Error message
			return nil, nil
		case 'X':
			// Termination
			return ParseTermination(w.Reader())
		default:
			return nil, fmt.Errorf("unknown message tag '%c'", start)
		}
	}
}
