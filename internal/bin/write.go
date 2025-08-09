package bin

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/util"
)

var endian = binary.BigEndian

const (
	strEnd       byte = 0x1f
	negativeSign byte = 0x2d
	positiveSign byte = 0x2b
)

func (b *BinDB) write(v byte) {
	b.buf.WriteByte(v)
}

func (b *BinDB) writeUint16(v uint16) {
	values := make([]byte, 2)
	endian.PutUint16(values, v)

	b.buf.Write(values)
}

func (b *BinDB) writeInt16(v int16) {
	if v < 0 {
		b.write(negativeSign)
		b.writeUint16(uint16(-v))
	} else {
		b.write(positiveSign)
		b.writeUint16(uint16(v))
	}
}

func (b *BinDB) writeString(f string) {
	b.buf.WriteString(f)
	b.buf.WriteByte(strEnd)
}

func (b *BinDB) writeBool(f bool) {
	if f {
		b.write(1)
	} else {
		b.write(0)
	}
}

func (b *BinDB) WriteStoreable(s db.Storeable) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	// kind
	b.write(kind.KindOf(s).Byte())
	// ID
	b.writeUint16(s.GetID().UInt16())

	// dump rest of the fields
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if field.Name == "ID" {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Uint8:
			b.write(byte(value.Uint()))
		case reflect.Uint16:
			b.writeUint16(uint16(value.Uint()))
		case reflect.Int16:
			b.writeInt16(int16(value.Int()))
		case reflect.String:
			b.writeString(value.String())
		case reflect.Bool:
			b.write(util.BoolToByte(value.Bool()))
		default:
			return fmt.Errorf("export: unsupported type %q", field.Type.Kind())
		}
	}

	return nil
}
