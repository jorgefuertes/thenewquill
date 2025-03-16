package db

import (
	"slices"
	"strconv"

	"thenewquill/internal/compiler/section"
	"thenewquill/internal/util"
)

type Register struct {
	Section section.Section
	Label   string
	Fields  []string
}

func NewRegister(section section.Section, label string, fields ...any) Register {
	r := Register{Section: section, Label: label, Fields: make([]string, 0)}

	for _, f := range fields {
		r.Fields = append(r.Fields, util.ValueToString(f))
	}

	return r
}

func (r *Register) GetString() string {
	if len(r.Fields) == 0 {
		return ""
	}

	str := r.Fields[0]
	r.Fields = slices.Delete(r.Fields, 0, 1)

	return str
}

func (r *Register) GetInt() int {
	i, err := strconv.Atoi(r.GetString())
	if err != nil {
		return 0
	}

	return i
}

func (r *Register) GetFloat() float64 {
	f, err := strconv.ParseFloat(r.GetString(), 64)
	if err != nil {
		return 0
	}

	return f
}

func (r *Register) GetBool() bool {
	b, err := strconv.ParseBool(r.GetString())
	if err != nil {
		return false
	}

	return b
}
