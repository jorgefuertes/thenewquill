package db

import (
	"fmt"

	"thenewquill/internal/compiler/section"
)

type Record struct {
	Section section.Section
	Label   string
	Fields  []any
}

func NewRecord(section section.Section, label string, fields ...any) Record {
	return Record{Section: section, Label: label, Fields: fields}
}

func (r *Record) Append(fields ...any) {
	r.Fields = append(r.Fields, fields...)
}

func (r Record) FieldAsMap(index int) map[string]any {
	v, ok := r.Fields[index].(map[string]any)
	if ok {
		return v
	}

	return map[string]any{}
}

func (r Record) FieldAsString(index int) string {
	v, ok := r.Fields[index].(string)
	if ok {
		return v
	}

	return fmt.Sprintf("%v", r.Fields[index])
}

func (r Record) FieldAsByte(index int) byte {
	v, ok := r.Fields[index].(byte)
	if ok {
		return v
	}

	return 0
}

func (r Record) FieldAsInt(index int) int {
	v, ok := r.Fields[index].(int)
	if ok {
		return v
	}

	return 0
}

func (r Record) FieldAsFloat(index int) float64 {
	v, ok := r.Fields[index].(float64)
	if ok {
		return v
	}

	return 0.0
}

func (r Record) FieldAsBool(index int) bool {
	b, ok := r.Fields[index].(bool)
	if ok {
		return b
	}

	return false
}
