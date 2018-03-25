package pgmsg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type readBuffer struct {
	*bufio.Reader
	buffer []byte
}

func newReadBuffer(r io.Reader) *readBuffer {
	// If we already have a read buffer, don't create a new one
	if buf, ok := r.(*readBuffer); ok {
		return buf
	}

	buf := &readBuffer{
		Reader: bufio.NewReader(r),
	}
	buf.Reader.Reset(r)
	return buf
}

func (b *readBuffer) PeekByte() (byte, error) {
	buf, err := b.Peek(1)
	if len(buf) == 1 {
		return buf[0], err
	}

	return '0', err
}

func (b *readBuffer) ReadInt() (int, error) {
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

func (b *readBuffer) ReadLength() (*readBuffer, error) {
	l, err := b.ReadInt()
	if err != nil {
		return nil, err
	}
	if l <= 0 {
		return nil, fmt.Errorf("unable to parse length from message")
	}

	// Length needs to account for the 4 bytes of the length value that have already been parsed
	l = l - 4
	if l == 0 {
		return nil, nil
	}

	buf := make([]byte, l)
	n, err := b.Read(buf)
	if err != nil {
		return nil, err
	}

	if l != n {
		return nil, fmt.Errorf("could not parse required bytes from message")
	}

	return newReadBuffer(bytes.NewReader(buf)), nil
}

func (b *readBuffer) ReadString() ([]byte, error) {
	str, err := b.ReadBytes(0)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return bytes.TrimRight(str, "\x00"), nil
}

func (b *readBuffer) ReadTag(t byte) error {
	tag, err := b.ReadByte()
	if err != nil {
		return err
	}
	if tag != t {
		return fmt.Errorf("invalid tag '%c' for message, must be '%c'", tag, t)
	}
	return nil
}

type writeBuffer struct {
	bytes []byte
}

func newWriteBuffer() *writeBuffer {
	return &writeBuffer{}
}

func (b *writeBuffer) Bytes() []byte {
	return b.bytes
}

func (b *writeBuffer) WriteInt(i int) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	b.bytes = append(b.bytes, buf...)
}

func (b *writeBuffer) WriteByte(c byte) {
	b.bytes = append(b.bytes, c)
}

func (b *writeBuffer) WriteString(buf []byte, null bool) {
	b.bytes = append(b.bytes, buf...)
	if null {
		b.bytes = append(b.bytes, '\x00')
	}
}

func (b *writeBuffer) PrependLength() {
	// Need to include the 4 bytes as part of the length
	l := len(b.bytes) + 4
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(l))
	b.bytes = append(buf, b.bytes...)
}

func (b *writeBuffer) PrependByte(c byte) {
	b.bytes = append([]byte{c}, b.bytes...)
}

func (b *writeBuffer) Wrap(t byte) {
	b.PrependLength()
	b.PrependByte(t)
}
