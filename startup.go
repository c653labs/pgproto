package pgproto

import (
	"bytes"
	"fmt"
	"io"
	"sort"
)

type StartupMessage struct {
	Options map[string][]byte
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
	p, err := buf.ReadInt()
	if err != nil {
		return nil, err
	}
	if p != ProtocolVersion {
		return nil, fmt.Errorf("unsupported protocol version")
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
	w.WriteInt(ProtocolVersion)

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
	str := fmt.Sprintf("StartupMessage<Protocol=%#v, Options<", ProtocolVersion)

	keys := make([]string, 0)
	for k, _ := range s.Options {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, k := range keys {
		if i > 0 {
			str += ", "
		}
		str += fmt.Sprintf("%s=%#v", k, string(s.Options[k]))
	}
	return str + ">>"
}
