package pgproto

import (
	"io"
)

type CloseType byte

const (
	CloseTypePreparedStatement CloseType = 'S'
	CloseTypePortal                      = 'P'
)

func (c CloseType) String() string {
	switch c {
	case CloseTypePreparedStatement:
		return "PreparedStatement"
	case CloseTypePortal:
		return "Portal"
	}
	return "Uknown"
}

type Close struct {
	ObjectType CloseType
	Name       []byte
}

func (c *Close) client() {}

func ParseClose(r io.Reader) (*Close, error) {
	b := newReadBuffer(r)

	err := b.ReadTag('C')
	if err != nil {
		return nil, err
	}

	c := &Close{}
	t, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	c.ObjectType = CloseType(t)
	c.Name, err = b.ReadString(stripNull)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Close) Encode() []byte {
	b := newWriteBuffer()
	b.WriteByte(byte(c.ObjectType))
	b.WriteString(c.Name, writeNull)
	b.Wrap('C')
	return b.Bytes()
}

func (c *Close) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "Close",
		"Payload": map[string]interface{}{
			"ObjectType": c.ObjectType,
			"Name":       c.Name,
		},
	}
}

func (c *Close) WriteTo(w io.Writer) (int64, error) { return writeTo(c, w) }
func (c *Close) String() string                     { return messageToString(c) }
