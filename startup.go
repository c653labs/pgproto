package pgmsg

import (
	"bytes"
	"io"
)

type StartupMessage struct {
	Protocol int
	Options  map[string][]byte
}

func ParseStartupMessage(r io.Reader) (*StartupMessage, error) {
	b := NewReadBuffer(r)

	// [int32 - length] [int32 - protocol] [[string]\0[string\0]]\0
	_, raw, err := b.ReadLength()
	if err != nil {
		return nil, err
	}

	// Replace the passed in buffer with one that is only scoped to the desired length we need
	b = NewReadBuffer(bytes.NewReader(raw))

	s := &StartupMessage{}

	// Parse protocol version
	s.Protocol, err = b.ReadInt()
	if err != nil {
		return nil, err
	}

	// Parse the key/value pairs
	for {
		key, err := b.ReadString()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		value, err := b.ReadString()
		if err != nil {
			return nil, err
		}

		s.Options[string(bytes.ToLower(key))] = value
	}

	return s, nil
}

func (s *StartupMessage) Encode() []byte {
	w := NewWriteBuffer()
	w.WriteInt(s.Protocol)
	for k, v := range s.Options {
		w.WriteString([]byte(k))
		w.WriteString(v)
	}
	w.PrependLength()

	return w.Bytes()
}
