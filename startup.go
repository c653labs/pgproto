package pgproto

import (
	"bytes"
	"fmt"
	"io"
	"sort"
)

type StartupMessage struct {
	Protocol int
	Options  map[string][]byte
}

func (s *StartupMessage) client() {}

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
		key, err := buf.ReadString(false)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		// This message ends in a single null terminator
		if bytes.Equal(key, []byte{'\x00'}) {
			break
		}
		// The key is [string] \0, we keep the \0 until now for the previous check
		key = bytes.TrimRight(key, "\x00")

		value, err := buf.ReadString(true)
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

func (s *StartupMessage) WriteTo(w io.Writer) (int64, error) { return writeTo(s, w) }

func (s *StartupMessage) String() string {
	return fmt.Sprintf("StartupMessage<Protocol=%#v, Options=%#v>", s.Protocol, s.Options)
}
