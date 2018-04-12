package pgproto

import (
	"bytes"
	"fmt"
	"io"
)

type Error struct {
	Severity []byte
	Text     []byte
	Code     []byte
	Message  []byte
	Position []byte
	File     []byte
	Line     []byte
	Routine  []byte
}

func (e *Error) server() {}

func ParseError(r io.Reader) (*Error, error) {
	b := newReadBuffer(r)

	// 'E'|'N' [int32 - length] ([char - key] [string - value] \0)+ \0
	tag, err := b.ReadByte()
	if err != nil {
		return nil, err
	}

	if tag != 'E' && tag != 'N' {
		return nil, fmt.Errorf("expected tag 'E' or 'N'")
	}

	b, err = b.ReadLength()
	if err != nil {
		return nil, err
	}

	e := &Error{}
	for {
		value, err := b.ReadString(false)
		if err != nil {
			return nil, err
		}

		// This message ends with a single null terminator
		if bytes.Equal(value, []byte{'\x00'}) {
			break
		}

		// Strip null terminator from the end
		value = bytes.TrimRight(value, "\x00")

		code := value[0]
		value = value[1:]
		switch code {
		case 'S':
			e.Severity = value
		case 'V':
			e.Text = value
		case 'C':
			e.Code = value
		case 'M':
			e.Message = value
		case 'P':
			e.Position = value
		case 'F':
			e.File = value
		case 'L':
			e.Line = value
		case 'R':
			e.Routine = value
		}
	}

	return e, nil
}

func (e *Error) Encode() []byte {
	return encodeError(e, 'E')
}

func (e *Error) WriteTo(w io.Writer) (int64, error) { return writeTo(e, w) }

func (e *Error) String() string {
	return errorString(e, "Error")
}

func encodeError(e *Error, tag byte) []byte {
	b := newWriteBuffer()

	// Severity
	b.WriteByte('S')
	b.WriteString(e.Severity, true)

	// Text
	b.WriteByte('V')
	b.WriteString(e.Text, true)

	// Code
	b.WriteByte('C')
	b.WriteString(e.Code, true)

	// Message
	b.WriteByte('M')
	b.WriteString(e.Message, true)

	// Position
	b.WriteByte('P')
	b.WriteString(e.Position, true)

	// Line
	b.WriteByte('L')
	b.WriteString(e.Line, true)

	// Routine
	b.WriteByte('R')
	b.WriteString(e.Routine, true)

	// Finalize
	b.WriteByte('\x00')
	b.Wrap(tag)
	return b.Bytes()
}

func errorString(e *Error, name string) string {
	return fmt.Sprintf(
		"%s<Severity=%#v, Text=%#v, Code=%#v, Message=%#v, Position=%#v, Line=%#v, Routine=%#v>",
		name,
		string(e.Severity),
		string(e.Text),
		string(e.Code),
		string(e.Message),
		string(e.Position),
		string(e.Line),
		string(e.Routine),
	)
}
