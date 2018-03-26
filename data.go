package pgproto

import (
	"fmt"
	"io"
)

type DataRow struct {
	Fields [][]byte
}

func ParseDataRow(r io.Reader) (*DataRow, error) {
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

	d := &DataRow{
		Fields: make([][]byte, c),
	}

	for i := 0; i < c; i++ {
		// [int32 - length] [string - data]
		l, err := b.ReadInt()
		if err != nil {
			return nil, err
		}
		d.Fields[i] = make([]byte, l)
		_, err = b.Read(d.Fields[i])
		if err != nil {
			return nil, err
		}
	}

	return d, nil
}

func (d *DataRow) Encode() []byte {
	b := newWriteBuffer()
	b.WriteInt16(len(d.Fields))
	for _, f := range d.Fields {
		b.WriteInt(len(f))
		b.WriteBytes(f)
	}
	b.Wrap('D')
	return b.Bytes()
}

func (d *DataRow) WriteTo(w io.Writer) (int, error) {
	return w.Write(d.Encode())
}

func (d *DataRow) String() string {
	str := "DataRow<"
	for i, f := range d.Fields {
		if i > 0 {
			str += ", "
		}
		str += fmt.Sprintf("%#v", string(f))
	}
	return str + ">"
}
