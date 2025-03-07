package db

import (
	"slices"
	"strconv"

	"thenewquill/internal/compiler/section"
)

type Register struct {
	Section section.Section
	Fields  []string
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
