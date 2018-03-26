package pgproto

import (
	"fmt"
	"io"
)

type RowFormat int

const (
	ROW_FORMAT_TEXT RowFormat = 0
)

func (f RowFormat) String() string {
	switch f {
	case ROW_FORMAT_TEXT:
		return "Text"
	}
	return "Unknown"
}

type RowField struct {
	ColumnName   []byte
	TableOID     int
	ColumnIndex  int // int16
	TypeOID      int
	ColumnLength int //int16
	TypeModifier int
	Format       RowFormat // int16
}

func (f RowField) String() string {
	return fmt.Sprintf(
		"%s<TableOID=%#v, ColumnIndex=%#v, TypeOID=%#v, ColumnLength=%#v, TypeModifier=%#v, Format=%s>",
		f.ColumnName, f.TableOID, f.ColumnIndex, f.TypeOID, f.ColumnLength, f.TypeModifier, f.Format,
	)
}

type RowDescription struct {
	Fields []RowField
}

func ParseRowDescription(r io.Reader) (*RowDescription, error) {
	b := newReadBuffer(r)

	// 'T' [int32 - length] [int32 - field count] ([string - column name]\0 [int32 - table oid] [int16 - column index] [int32 - type oid] [int16 - column length] [int32 - type modifier] [int16 - format])
	err := b.ReadTag('T')
	if err != nil {
		return nil, err
	}

	// Length - int
	b, err = b.ReadLength()
	if err != nil {
		return nil, err
	}

	// Field count - int16
	c, err := b.ReadInt16()
	if err != nil {
		return nil, err
	}

	rd := &RowDescription{
		Fields: make([]RowField, c),
	}
	for i := 0; i < c; i++ {
		// Column Name - string
		rd.Fields[i].ColumnName, err = b.ReadString(true)
		if err != nil {
			return nil, err
		}

		// Table OID - int
		rd.Fields[i].TableOID, err = b.ReadInt()
		if err != nil {
			return nil, err
		}

		// Column Index - int16
		rd.Fields[i].ColumnIndex, err = b.ReadInt16()
		if err != nil {
			return nil, err
		}

		// Type OID - int
		rd.Fields[i].TypeOID, err = b.ReadInt()
		if err != nil {
			return nil, err
		}

		// Column Length - int16
		rd.Fields[i].ColumnLength, err = b.ReadInt16()
		if err != nil {
			return nil, err
		}

		// Type Modifier - int
		rd.Fields[i].TypeModifier, err = b.ReadInt()
		if err != nil {
			return nil, err
		}

		// Format - int16
		format, err := b.ReadInt16()
		rd.Fields[i].Format = RowFormat(format)
		if err != nil {
			return nil, err
		}
	}

	return rd, nil
}

func (r *RowDescription) Encode() []byte {
	b := newWriteBuffer()
	// Field count - int16
	b.WriteInt16(len(r.Fields))
	for _, f := range r.Fields {
		// Column Name - string
		b.WriteString(f.ColumnName, true)

		// Table OID - int
		b.WriteInt(f.TableOID)

		// Column Index - int16
		b.WriteInt16(f.ColumnIndex)

		// Type OID - int
		b.WriteInt(f.TypeOID)

		// Column Length - int16
		b.WriteInt16(f.ColumnLength)

		// Type Modifier - int
		b.WriteInt(f.TypeModifier)

		// Format - int16
		b.WriteInt16(int(f.Format))
	}

	b.Wrap('T')
	return b.Bytes()
}

func (r *RowDescription) WriteTo(w io.Writer) (int64, error) { return writeTo(r, w) }

func (r *RowDescription) String() string {
	str := "RowDescription<"
	for i, f := range r.Fields {
		if i > 0 {
			str += ", "
		}
		str += f.String()
	}
	return str + ">"
}
