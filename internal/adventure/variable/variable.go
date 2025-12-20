package variable

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/pkg/log"
)

const (
	FalseValue = "0"
	TrueValue  = "1"
)

var TrueValuesRg = regexp.MustCompile(`^(?i)(1|t|true|yes|y|s|si|sí|on)$`)

var (
	RegexpInt   = regexp.MustCompile(`^\d+$`)
	RegexpFloat = regexp.MustCompile(`^\d+\.\d+$`)
)

type Variable struct {
	ID      uint32
	LabelID uint32
	Value   string
}

var _ adapter.Storeable = &Variable{}

func New(id, labelID uint32, value any) *Variable {
	v := &Variable{ID: id, LabelID: labelID}
	v.Set(value)

	return v
}

func (v *Variable) SetID(id uint32) {
	v.ID = id
}

func (v Variable) GetID() uint32 {
	return v.ID
}

func (v *Variable) SetLabelID(id uint32) {
	v.ID = id
}

func (v Variable) GetLabelID() uint32 {
	return v.ID
}

func (v *Variable) Set(value any) {
	switch val := value.(type) {
	case bool:
		if val {
			v.Value = TrueValue
		} else {
			v.Value = FalseValue
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		v.Value = fmt.Sprintf("%d", val)
	case float32, float64:
		v.Value = fmt.Sprintf("%.4f", val)
	case string:
		v.Value = val
	default:
		log.Warning("storing unknown value type %+v into var %d", val, v.ID)
		v.Value = fmt.Sprint(value)
	}
}

func (v Variable) String() string {
	return v.Value
}

func (v Variable) Int() int {
	if !v.IsInt() {
		return 0
	}

	i, err := strconv.Atoi(v.Value)
	if err != nil {
		log.Warning("error while parsing var to int: %v", err)

		return 0
	}

	return i
}

func (v Variable) Float() float64 {
	if !v.IsNum() {
		return 0
	}

	f, err := strconv.ParseFloat(v.Value, 64)
	if err != nil {
		log.Warning("error while parsing var to float: %v", err)

		return 0
	}

	return f
}

func (v Variable) Bool() bool {
	if TrueValuesRg.MatchString(v.Value) {
		return true
	}

	return v.Float() >= 1
}

func (v Variable) IsNum() bool {
	return v.IsInt() || v.IsFloat()
}

func (v Variable) IsInt() bool {
	return RegexpInt.MatchString(v.Value)
}

func (v Variable) IsFloat() bool {
	return RegexpFloat.MatchString(v.Value)
}

func (v Variable) IsTrue() bool {
	return v.Bool()
}

func (v Variable) IsFalse() bool {
	return !v.Bool()
}

func (v *Variable) IsEqual(other *Variable) bool {
	return v.Value == other.Value
}
