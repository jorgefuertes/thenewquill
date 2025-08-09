package variable

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

type Variable struct {
	ID    db.ID
	Value any
}

var _ db.Storeable = Variable{}

func New(val any) Variable {
	return Variable{ID: db.UndefinedLabel.ID, Value: val}
}

func (v Variable) SetID(id db.ID) db.Storeable {
	v.ID = id

	return v
}

func (v Variable) GetID() db.ID {
	return v.ID
}

func (v Variable) Int() int {
	switch v.Value.(type) {
	case int:
		return v.Value.(int)
	case byte:
		return int(v.Value.(byte))
	case float64:
		return int(v.Value.(float64))
	case bool:
		if v.Value.(bool) {
			return 1
		}

		return 0
	case string:
		i, _ := strconv.Atoi(v.Value.(string))

		return i
	default:
		return 0
	}
}

func (v Variable) Float() float64 {
	switch v.Value.(type) {
	case int:
		return float64(v.Value.(int))
	case float64:
		return v.Value.(float64)
	case bool:
		if v.Value.(bool) {
			return 1
		}

		return 0
	case string:
		f, _ := strconv.ParseFloat(v.Value.(string), 64)

		return f
	default:
		return 0
	}
}

func (v Variable) String() string {
	switch v.Value.(type) {
	case string:
		return v.Value.(string)
	case bool:
		return strconv.FormatBool(v.Value.(bool))
	case float64:
		return fmt.Sprintf("%.2f", v.Value.(float64))
	case int:
		return strconv.Itoa(v.Value.(int))
	case byte:
		return string(v.Value.(byte))
	case []byte:
		return string(v.Value.([]byte))
	default:
		return fmt.Sprint(v.Value)
	}
}

func (v Variable) Bool() bool {
	switch v.Value.(type) {
	case bool:
		return v.Value.(bool)
	case byte:
		return v.Value.(byte) >= 1
	case float64:
		return v.Value.(float64) >= 1
	case int:
		return v.Value.(int) >= 1
	case string:
		return regexp.MustCompile(`(?i)^([sy1t]{1}|si|sí|yes|true|on)$`).MatchString(v.Value.(string))
	default:
		return false
	}
}

func (v Variable) Byte() byte {
	switch v.Value.(type) {
	case byte:
		return v.Value.(byte)
	case int:
		return byte(v.Value.(int))
	case bool:
		if v.Value.(bool) {
			return 1
		}

		return 0
	case string:
		if v.Value.(string) == "true" {
			return 1
		}

		if v.Value.(string) == "false" {
			return 0
		}

		return byte(v.Value.(string)[0])
	default:
		return 0
	}
}
