package pgproto

import (
	"fmt"
	"io"
)

type ReadyStatus int

const (
	READY_IDLE ReadyStatus = 73
)

func (r ReadyStatus) String() string {
	switch r {
	case READY_IDLE:
		return "Idle"
	}
	return "Unknown"
}

type ReadyForQuery struct {
	Status ReadyStatus
}

func ParseReadyForQuery(r io.Reader) (*ReadyForQuery, error) {
	b := newReadBuffer(r)

	// 'Z' [int32 - length] [byte - status]
	err := b.ReadTag('Z')
	if err != nil {
		return nil, err
	}

	l, err := b.ReadInt()
	if err != nil {
		return nil, err
	}
	if l != 5 {
		return nil, fmt.Errorf("unexpected message length")
	}

	i, err := b.ReadByte()
	if err != nil {
		return nil, err
	}

	return &ReadyForQuery{
		Status: ReadyStatus(i),
	}, nil
}

func (r *ReadyForQuery) Encode() []byte {
	b := newWriteBuffer()
	b.WriteByte(byte(r.Status))
	b.Wrap('Z')
	return b.Bytes()
}

func (r *ReadyForQuery) WriteTo(w io.Writer) (int, error) {
	return w.Write(r.Encode())
}

func (r *ReadyForQuery) String() string {
	return fmt.Sprintf("ReadyForQuery<Status=%s>", r.Status)
}
