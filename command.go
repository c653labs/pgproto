package pgproto

import (
	"io"
)

type CommandCompletion struct {
	Tag []byte
}

func (c *CommandCompletion) server() {}

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

func (c *CommandCompletion) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "CommandCompletion",
		"Payload": map[string]string{
			"Tag": string(c.Tag),
		},
	}
}

func (c *CommandCompletion) WriteTo(w io.Writer) (int64, error) { return writeTo(c, w) }
func (c *CommandCompletion) String() string                     { return messageToString(c) }
