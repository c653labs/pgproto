package pgproto

import (
	"bytes"
	"fmt"
	"io"
)

// Message is the main interface for all PostgreSQL messages
type Message interface {
	Encode() []byte
	AsMap() map[string]interface{}
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

	msg, err := readRawMessage(r)
	if err != nil {
		return nil, err
	}
	start := msg[0]
	msgReader := bytes.NewReader(msg)

	// Startup message:
	//   [int32 - length] [int32 - protocol] [[string]\0[string]\0] \0
	// Regular message
	//   [char - tag] [int32 - length] [payload]
	switch start {
	case '\x00':
		// TODO: We need to handle this case better, it might not always start with \x00
		//       We could just make calling `ParseStartupMessage` explicit
		return ParseStartupMessage(msgReader)
	case 'p':
		// Password message
		return ParsePasswordMessage(msgReader)
	case 'Q':
		// Simple query
		return ParseSimpleQuery(msgReader)
	case 't':
		// Parameter description
		return ParseParameterDescription(msgReader)
	case 'B':
		// Binary parameters
		return ParseBinaryParameters(msgReader)
	case 'P':
		// Parse
		return ParseParse(msgReader)
	case 'E':
		// Execute
		return ParseExecute(msgReader)
	case 'H':
		// Flush
		return ParseFlush(msgReader)
	case 'S':
		// Sync
		return ParseSync(msgReader)
	case 'C':
		// Close
		return ParseClose(msgReader)
	case 'D':
		// Describe
		return ParseDescribe(msgReader)
	case 'X':
		// Termination
		return ParseTermination(msgReader)
	default:
		return nil, fmt.Errorf("unknown message tag '%c'", start)
	}
}

// ParseServerMessage will read the next ServerMessage from the provided io.Reader
func ParseServerMessage(r io.Reader) (ServerMessage, error) {

	msg, err := readRawMessage(r)
	if err != nil {
		return nil, err
	}
	start := msg[0]
	msgReader := bytes.NewReader(msg)
	// Message
	//   [char - tag] [int32 - length] [payload]
	switch start {
	case 'R':
		// Authentication request
		return ParseAuthenticationRequest(msgReader)
	case 'S':
		// Parameter status
		return ParseParameterStatus(msgReader)
	case 'K':
		// Backend key data
		return ParseBackendKeyData(msgReader)
	case 'Z':
		// Ready for query
		return ParseReadyForQuery(msgReader)
	case 'C':
		// Command completion
		return ParseCommandCompletion(msgReader)
	case 'T':
		// Row description
		return ParseRowDescription(msgReader)
	case 't':
		// Parameter description
		return nil, fmt.Errorf("unhandled message tag %#v", start)
	case 'D':
		// Data row
		return ParseDataRow(msgReader)
	case 'I':
		// Empty query response
		return ParseEmptyQueryResponse(msgReader)
	case '1':
		// Parse complete
		return ParseParseComplete(msgReader)
	case '2':
		// Bind complete
		return ParseBindComplete(msgReader)
	case '3':
		// Close complete
		return ParseCloseComplete(msgReader)
	case 'W':
		// Copy both response
		return ParseCopyBothResponse(msgReader)
	case 'd':
		// Copy data
		return ParseCopyData(msgReader)
	case 'G':
		// Copy in response
		return ParseCopyInResponse(msgReader)
	case 'H':
		// Copy out response
		return ParseCopyOutResponse(msgReader)
	case 'V':
		// Function call response
		return nil, fmt.Errorf("unhandled message tag %#v", start)
	case 'n':
		// No data
		return ParseNoData(msgReader)
	case 'N':
		// Notice response
		return ParseNoticeResponse(msgReader)
	case 'A':
		// Notification response
		return ParseNotification(msgReader)
	case 'E':
		// Error message
		return ParseError(msgReader)
	default:
		return nil, fmt.Errorf("unknown message tag '%c'", start)
	}
}

func readRawMessage(r io.Reader) (rawmsg []byte, err error) {

	startByte, err := ReadNBytes(r, 1)
	if err != nil {
		return nil, err
	}
	rawPkgLen, err := ReadNBytes(r, 4)
	if err != nil {
		return nil, err
	}
	pkgLen := bytesToInt(rawPkgLen)
	payload, err := ReadNBytes(r, pkgLen-4)
	if err != nil {
		return nil, err
	}
	rawMessage := make([]byte, pkgLen+1)
	copy(rawMessage, startByte)
	copy(rawMessage[1:], rawPkgLen)
	copy(rawMessage[5:], payload)
	return rawMessage, nil

}
