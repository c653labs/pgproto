package pgproto

import (
	"fmt"
	"io"
)

type CommandCompletion struct {
	Tag []byte
}

func ParseCommandCompletion(r io.Reader) (*CommandCompletion, error) {
	b := newReadBuffer(r)

	// 'C' [int32 - length] [tag] \0
	err := b.ReadTag('C')
	if err != nil {
		return nil, err
	}

	b, err = b.ReadLength()
	if err != nil {
		return nil, err
	}

	c := &CommandCompletion{}
	c.Tag, err = b.ReadString(true)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CommandCompletion) Encode() []byte {
	b := newWriteBuffer()
	b.WriteString(c.Tag, true)
	b.Wrap('C')
	return b.Bytes()
}

func (c *CommandCompletion) WriteTo(w io.Writer) (int, error) {
	return w.Write(c.Encode())
}

func (c *CommandCompletion) String() string {
	return fmt.Sprintf("CommandCompletion<Tag=%#v>", string(c.Tag))
}
