package bin

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/pkg/log"
)

var endian = binary.BigEndian

const (
	strEnd       byte = 0x1f
	negativeSign byte = 0x2d
	positiveSign byte = 0x2b
	sliceBegin   byte = 0x1d
	sliceEnd     byte = 0x1e
	sliceSep     byte = 0x1c
)

func (b *BinDB) write(v byte) {
	if err := b.buf.WriteByte(v); err != nil {
		log.Fatal("Unexpected error writing to buffer: %s", err)
	}
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

func (b *BinDB) writeUint32(v uint32) {
	values := make([]byte, 2)
	endian.PutUint32(values, v)

	b.buf.Write(values)
}

func (b *BinDB) writeInt32(v int32) {
	if v < 0 {
		b.write(negativeSign)
		b.writeUint32(uint32(-v))
	} else {
		b.write(positiveSign)
		b.writeUint32(uint32(v))
	}
}

func (b *BinDB) writeUint64(v uint64) {
	values := make([]byte, 2)
	endian.PutUint64(values, v)

	b.buf.Write(values)
}

func (b *BinDB) writeInt64(v int64) {
	if v < 0 {
		b.write(negativeSign)
		b.writeUint64(uint64(-v))
	} else {
		b.write(positiveSign)
		b.writeUint64(uint64(v))
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

		if err := b.writeValue(value); err != nil {
			return fmt.Errorf("error writing storeable %q field %q: %s", kind.KindOf(s), field.Name, err)
		}
	}

	return nil
}

func (b *BinDB) writeValue(v reflect.Value) error {
	switch v.Type().Kind() {
	case reflect.Uint8:
		b.write(byte(v.Uint()))
	case reflect.Uint16:
		b.writeUint16(uint16(v.Uint()))
	case reflect.Int16:
		b.writeInt16(int16(v.Int()))
	case reflect.Uint32:
		b.writeUint32(uint32(v.Uint()))
	case reflect.Int32:
		b.writeInt32(int32(v.Int()))
	case reflect.Uint64:
		b.writeUint64(uint64(v.Uint()))
	case reflect.Int64:
		b.writeInt64(int64(v.Int()))
	case reflect.String:
		b.writeString(v.String())
	case reflect.Bool:
		b.writeBool(v.Bool())
	case reflect.Slice:
		return b.writeSlice(v)
	default:
		return fmt.Errorf("writeValue: unsupported type %q", v.Type().Kind())
	}

	return nil
}

func (b *BinDB) writeSlice(v reflect.Value) error {
	b.write(sliceBegin)

	if v.Type().Elem().Kind() == reflect.String {
		b.writeString(strings.Join(v.Interface().([]string), ","))
	} else {
		return fmt.Errorf("writeSlice: unsupported type %q", v.Type().Elem().Kind())
	}

	b.write(sliceEnd)

	return nil
}
