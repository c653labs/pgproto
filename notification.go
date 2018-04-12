package pgproto

import (
	"fmt"
	"io"
)

type Notification struct {
	PID       int
	Condition []byte
}

func (n *Notification) server() {}

func ParseNotification(r io.Reader) (*Notification, error) {
	buf := newReadBuffer(r)

	// 'A' [int32 - length] [int32 - pid] [string - condition] \0
	err := buf.ReadTag('A')
	if err != nil {
		return nil, err
	}

	buf, err = buf.ReadLength()
	if err != nil {
		return nil, err
	}

	pid, err := buf.ReadInt()
	if err != nil {
		return nil, err
	}

	condition, err := buf.ReadString(true)
	if err != nil {
		return nil, err
	}
	return &Notification{
		PID:       pid,
		Condition: condition,
	}, nil
}

func (n *Notification) Encode() []byte {
	// 'A' [int32 - length] [int32 - pid] [string - condition] \0
	buf := newWriteBuffer()
	buf.WriteInt(n.PID)
	buf.WriteString(n.Condition, true)
	buf.Wrap('N')
	return buf.Bytes()
}

func (n *Notification) WriteTo(w io.Writer) (int64, error) { return writeTo(n, w) }

func (n *Notification) String() string {
	return fmt.Sprintf("Notification<PID=%#v, Condition=%#v>", n.PID, string(n.Condition))
}
