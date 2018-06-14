package pgproto

import (
	"bytes"
	"fmt"
	"io"
)

// RawMessage represents the common message format on the wire [char tag] [int32 length] [bytes payload]
type RawMessage struct {
	Type    byte
	Length  int32
	Payload []byte
}

func ParseRawMessage(r io.Reader) (*RawMessage, error) {
	buf := newReadBuffer(r)

	// [char tag] [int32 length] [payload]
	tag, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}

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

	return &RawMessage{
		Type:    tag,
		Length:  int32(l),
		Payload: b,
	}, nil
}

func (r *RawMessage) ValidateTag(b byte) error {
	if r.Type != b {
		return fmt.Errorf("invalid tag '%c' for message, must be '%c'", r.Type, b)
	}
	return nil
}

func (r *RawMessage) PayloadReader() *readBuffer {
	return newReadBuffer(bytes.NewReader(r.Payload))
}

func (r *RawMessage) writeBuffer() *writeBuffer {
	w := newWriteBuffer()
	w.WriteByte(r.Type)
	// Add 4 for the size of the byte we are writing
	w.WriteInt(len(r.Payload) + 4)
	w.WriteBytes(r.Payload)
	return w
}

func (r *RawMessage) Reader() io.Reader {
	w := r.writeBuffer()
	return w.Reader()
}

func (r *RawMessage) Encode() []byte {
	w := r.writeBuffer()
	return w.Bytes()
}

func (r *RawMessage) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Tag":     string(r.Type),
		"Payload": r.Payload,
	}
}

func (r *RawMessage) String() string { return messageToString(r) }
