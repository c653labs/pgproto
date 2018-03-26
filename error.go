package pgmsg

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

func ParseError(r io.Reader) (*Error, error) {
	b := newReadBuffer(r)

	// 'E' [int32 - length] ([char - key] [string - value] \0)+ \0
	err := b.ReadTag('E')
	if err != nil {
		return nil, err
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
	b.Wrap('E')
	return b.Bytes()
}

func (e *Error) WriteTo(w io.Writer) (int, error) {
	return w.Write(e.Encode())
}

func (e *Error) String() string {
	return fmt.Sprintf(
		"Error<Severity=%#v, Text=%#v, Code=%#v, Message=%#v, Position=%#v, Line=%#v, Routine=%#v>",
		string(e.Severity),
		string(e.Text),
		string(e.Code),
		string(e.Message),
		string(e.Position),
		string(e.Line),
		string(e.Routine),
	)
}
