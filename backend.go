package pgproto

import (
	"fmt"
	"io"
)

type BackendKeyData struct {
	PID int
	Key int
}

func ParseBackendKeyData(r io.Reader) (*BackendKeyData, error) {
	buf := newReadBuffer(r)

	// 'K' [int32 - length] [int32 - pid] [in32 - key]
	err := buf.ReadTag('K')
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

	key, err := buf.ReadInt()
	if err != nil {
		return nil, err
	}
	return &BackendKeyData{
		PID: pid,
		Key: key,
	}, nil
}

func (b *BackendKeyData) Encode() []byte {
	buf := newWriteBuffer()
	buf.WriteInt(b.PID)
	buf.WriteInt(b.Key)
	buf.Wrap('K')
	return buf.Bytes()
}

func (b *BackendKeyData) WriteTo(w io.Writer) (int64, error) { return writeTo(b, w) }

func (b *BackendKeyData) String() string {
	return fmt.Sprintf("BackendKeyData<PID=%#v, Key=%#v>", b.PID, b.Key)
}
