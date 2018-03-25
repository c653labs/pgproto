package pgmsg

import (
	"bytes"
	"io"
	"sort"
)

type StartupMessage struct {
	Protocol int
	Options  map[string][]byte
}

func ParseStartupMessage(r io.Reader) (*StartupMessage, error) {
	b := newReadBuffer(r)

	// [int32 - length] [int32 - protocol] [[string]\0[string\0]]\0
	buf, err := b.ReadLength()
	if err != nil {
		return nil, err
	}

	s := &StartupMessage{
		Options: make(map[string][]byte),
	}

	// Parse protocol version
	s.Protocol, err = buf.ReadInt()
	if err != nil {
		return nil, err
	}

	// Parse the key/value pairs
	for {
		// Check if the next byte is a null terminator
		// DEV: This message ends in a null terminator
		n, err := buf.PeekByte()
		if n == '\x00' {
			break
		} else if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		key, err := buf.ReadString()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		value, err := buf.ReadString()
		if err != nil {
			return nil, err
		}

		s.Options[string(bytes.ToLower(key))] = value
	}

	return s, nil
}

func (s *StartupMessage) Encode() []byte {
	w := newWriteBuffer()
	w.WriteInt(s.Protocol)

	// Encode the options in sorted order
	keys := []string{}
	for k := range s.Options {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := s.Options[k]
		w.WriteString([]byte(k), true)
		w.WriteString(v, true)
	}
	w.WriteByte('\x00')
	w.PrependLength()

	return w.Bytes()
}

func (s *StartupMessage) WriteTo(w io.Writer) (int, error) {
	return w.Write(s.Encode())
}
