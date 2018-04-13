package pgproto

import (
	"io"
)

type BinaryParameters struct {
	Fields [][]byte
}

func (p *BinaryParameters) client() {}

func ParseBinaryParameters(r io.Reader) (*BinaryParameters, error) {
	b := newReadBuffer(r)

	// 'D' [int32 - length] [int16 - field count] ([int32 - length] [string - data])+
	err := b.ReadTag('D')
	if err != nil {
		return nil, err
	}

	b, err = b.ReadLength()
	if err != nil {
		return nil, err
	}

	// Field count - int16
	c, err := b.ReadInt16()
	if err != nil {
		return nil, err
	}

	p := &BinaryParameters{
		Fields: make([][]byte, c),
	}

	for i := 0; i < c; i++ {
		// [int32 - length] [string - data]
		l, err := b.ReadInt()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if l == -1 {
			p.Fields[i] = nil
		} else {
			p.Fields[i] = make([]byte, l)
			_, err = b.Read(p.Fields[i])
			if err != nil {
				return nil, err
			}
		}
	}

	return p, nil
}

func (p *BinaryParameters) Encode() []byte {
	b := newWriteBuffer()
	b.WriteInt16(len(p.Fields))
	for _, f := range p.Fields {
		b.WriteInt(len(f))
		b.WriteBytes(f)
	}
	b.Wrap('D')
	return b.Bytes()
}

func (p *BinaryParameters) AsMap() map[string]interface{} {
	f := make([]string, len(p.Fields))
	for k, v := range p.Fields {
		f[k] = string(v)
	}
	return map[string]interface{}{
		"Type": "BinaryParameters",
		"Payload": map[string]interface{}{
			"Fields": f,
		},
	}
}

func (p *BinaryParameters) WriteTo(w io.Writer) (int64, error) { return writeTo(p, w) }
func (p *BinaryParameters) String() string                     { return messageToString(p) }
