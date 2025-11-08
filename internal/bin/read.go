package bin

import (
	"errors"
	"reflect"

	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/item"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	"github.com/jorgefuertes/thenewquill/internal/adventure/message"
	"github.com/jorgefuertes/thenewquill/internal/adventure/variable"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
)

var (
	ErrEOF = errors.New("end of file")

	ErrNextKindIsNone    = errors.New("kind is none")
	ErrNextKindIsUnknown = errors.New("unsupported kind")
)

func (b *BinDB) readKind() (kind.Kind, error) {
	if b.buf.Len() == 0 {
		return 0, ErrEOF
	}

	bKind, err := b.read()
	if err != nil {
		return kind.None, err
	}

	return kind.KindFromByte(bKind), nil
}

func (b *BinDB) readID() (db.ID, error) {
	id, err := b.readUint16()
	if err != nil {
		return 0, err
	}

	return db.ID(id), nil
}

func (b *BinDB) read() (byte, error) {
	return b.buf.ReadByte()
}

func (b *BinDB) readUint16() (uint16, error) {
	values := make([]byte, 2)
	if _, err := b.buf.Read(values); err != nil {
		return 0, err
	}

	return endian.Uint16(values), nil
}

func (b *BinDB) readInt16() (int16, error) {
	sign, err := b.buf.ReadByte()
	if err != nil {
		return 0, err
	}

	values := make([]byte, 2)
	_, err = b.buf.Read(values)
	if err != nil {
		return 0, err
	}

	n := endian.Uint16(values)

	if sign == negativeSign {
		return -int16(n), nil
	}

	return int16(n), nil
}

func (b *BinDB) readString() (string, error) {
	var s []byte

	for {
		c, err := b.buf.ReadByte()
		if err != nil {
			return "", nil
		}

		if c == strEnd {
			break
		}

		s = append(s, c)
	}

	return string(s), nil
}

func (b *BinDB) readBool() (bool, error) {
	c, err := b.buf.ReadByte()
	if err != nil {
		return false, err
	}

	return c == 1, nil
}

func (b *BinDB) ReadStoreable() (db.Storeable, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	k, err := b.readKind()
	if err != nil {
		return nil, err
	}

	id, err := b.readID()
	if err != nil {
		return nil, err
	}

	s, err := getStoreable(k, id)
	if err != nil {
		return nil, err
	}

	// read the rest of the fields
	v := reflect.ValueOf(s).Elem() // Get the underlying value
	t := v.Type()                  // Get the type from the value

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if field.Name == "ID" {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Uint8:
			n, err := b.read()
			if err != nil {
				return nil, err
			}

			value.SetUint(uint64(n))
		case reflect.Uint16:
			n, err := b.readUint16()
			if err != nil {
				return nil, err
			}

			value.SetUint(uint64(n))
		case reflect.Int16:
			n, err := b.readInt16()
			if err != nil {
				return nil, err
			}

			value.SetInt(int64(n))
		case reflect.String:
			str, err := b.readString()
			if err != nil {
				return nil, err
			}

			value.SetString(str)
		case reflect.Bool:
			boolValue, err := b.readBool()
			if err != nil {
				return nil, err
			}

			value.SetBool(boolValue)
		}
	}

	return reflect.ValueOf(s).Elem().Interface().(db.Storeable), nil
}

func getStoreable(k kind.Kind, id db.ID) (db.Storeable, error) {
	switch k {
	case kind.Label:
		return &db.Label{ID: id}, nil
	case kind.Param:
		return &config.Param{ID: id}, nil
	case kind.Variable:
		return &variable.Variable{ID: id}, nil
	case kind.Word:
		return &word.Word{ID: id}, nil
	case kind.Message:
		return &message.Message{ID: id}, nil
	case kind.Item:
		return &item.Item{ID: id}, nil
	case kind.Location:
		return &location.Location{ID: id}, nil
	case kind.Character:
		return &character.Character{ID: id}, nil
	case kind.None:
		return nil, ErrNextKindIsNone
	default:
		return nil, ErrNextKindIsUnknown
	}
}
