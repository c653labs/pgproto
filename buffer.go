package pgmsg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type ReadBuffer struct {
	*bufio.Reader
}

func NewReadBuffer(r io.Reader) *ReadBuffer {
	buf := &ReadBuffer{
		Reader: bufio.NewReader(r),
	}
	buf.Reader.Reset(r)
	return buf
}

func (b *ReadBuffer) PeekByte() (byte, error) {
	buf, err := b.Peek(1)
	if len(buf) == 1 {
		return buf[0], err
	}

	return '0', err
}

func (b *ReadBuffer) ReadInt() (int, error) {
	buf := make([]byte, 4)
	n, err := b.Read(buf)
	if err != nil {
		return 0, err
	}
	if n != 4 {
		return 0, io.EOF
	}

	return int(int32(binary.BigEndian.Uint32(buf))), nil
}

func (b *ReadBuffer) ReadLength() (int, []byte, error) {
	l, err := b.ReadInt()
	if err != nil {
		return 0, nil, err
	}
	if l <= 0 {
		return 0, nil, fmt.Errorf("unable to parse length from message")
	}

	// Length needs to account for the 4 bytes of the length value that have already been parsed
	l = l - 4
	if l == 0 {
		return l, []byte{}, nil
	}

	buf := make([]byte, l)
	n, err := b.Read(buf)
	if err != nil {
		return 0, nil, err
	}

	if l != n {
		return 0, nil, fmt.Errorf("could not parse required bytes from message")
	}

	return l, buf, nil
}

func (b *ReadBuffer) ReadString() ([]byte, error) {
	str, err := b.ReadBytes(0)
	if err != nil {
		return nil, err
	}

	return bytes.TrimRight(str, "\x00"), nil
}

type WriteBuffer struct {
	bytes []byte
}

func NewWriteBuffer() *WriteBuffer {
	return &WriteBuffer{}
}

func (b *WriteBuffer) Bytes() []byte {
	return b.bytes
}

func (b *WriteBuffer) WriteInt(i int) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	b.bytes = append(b.bytes, buf...)
}

func (b *WriteBuffer) WriteByte(c byte) {
	b.bytes = append([]byte{c}, b.bytes...)
}

func (b *WriteBuffer) WriteString(buf []byte) {
	b.bytes = append(b.bytes, buf...)
	b.bytes = append(b.bytes, '\x00')
}

func (b *WriteBuffer) PrependLength() {
	// Need to include the 4 bytes as part of the length
	l := len(b.bytes) + 4
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(l))
	b.bytes = append(buf, b.bytes...)
}

func (b *WriteBuffer) PrependByte(c byte) {
	b.bytes = append([]byte{c}, b.bytes...)
}
