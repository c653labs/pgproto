package pgproto

import (
	"fmt"
	"io"
)

// Message is the main interface for all PostgreSQL messages
type Message interface {
	Encode() []byte
	WriteTo(w io.Writer) (int64, error)
	String() string
}

// ClientMessage is an interface describing all client side PostgreSQL messages (messages sent to the server)
type ClientMessage interface {
	Message
	client()
}

// ServerMessage is an interface describing all server side PostgreSQL messages (messages sent to the client)
type ServerMessage interface {
	Message
	server()
}

// ParseClientMessage will read the next ClientMessage from the provided io.Reader
func ParseClientMessage(r io.Reader) (ClientMessage, error) {
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
		case 'Q':
			// Simple query
			return ParseSimpleQuery(w.Reader())
		case 't':
			// Parameter description
			return nil, nil
		case 'B':
			// Binary parameters
			return nil, nil
		case 'P':
			// Parse
			return nil, nil
		case 'E':
			// Execute
			return nil, nil
		case 'H':
			// Flush
			return nil, nil
		case 'S':
			// Sync
			return nil, nil
		case 'C':
			// Close
			return nil, nil
		case 'D':
			// Describe
			return nil, nil
		case 'X':
			// Termination
			return ParseTermination(w.Reader())
		default:
			return nil, fmt.Errorf("unknown message tag '%c'", start)
		}
	}
}

// ParseServerMessage will read the next ServerMessage from the provided io.Reader
func ParseServerMessage(r io.Reader) (ServerMessage, error) {
	// Create a buffer
	buf := newReadBuffer(r)

	// Look at the first byte to determine the type of message we have
	start, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}

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
	case 'C':
		// Command completion
		return ParseCommandCompletion(w.Reader())
	case 'T':
		// Row description
		return ParseRowDescription(w.Reader())
	case 't':
		// Parameter description
		return nil, nil
	case 'D':
		// Data row
		return ParseDataRow(w.Reader())
	case 'I':
		// Empty query response
		return ParseEmptyQueryResponse(w.Reader())
	case 'B':
		// Bind
		return nil, nil
	case '2':
		// Bind complete
		return nil, nil
	case '3':
		// Close complete
		return nil, nil
	case 'W':
		// Copy both response
		return nil, nil
	case 'd':
		// Copy data
		return nil, nil
	case 'G':
		// Copy in response
		return nil, nil
	case 'H':
		// Copy out response
		return nil, nil
	case 'V':
		// Function call response
		return nil, nil
	case 'n':
		// No data
		return nil, nil
	case 'N':
		// Notice response
		return nil, nil
	case 'A':
		// Notification response
		return nil, nil
	case 'P':
		// Parse complete
		return nil, nil
	case 'E':
		// Error message
		return ParseError(w.Reader())
	default:
		return nil, fmt.Errorf("unknown message tag '%c'", start)
	}
}
