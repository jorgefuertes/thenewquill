package database

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type condition int

const (
	Equal condition = iota
	NotEqual
	Contains
	NotContains
)

type filter struct {
	condition condition
	field     string
	value     any
}

func Filter(field string, condition condition, value any) filter {
	return filter{condition, field, value}
}

func FilterByKind(k kind.Kind) filter {
	return filter{Equal, "kind", k}
}

func FilterByID(id primitive.ID) filter {
	return filter{Equal, "id", id}
}

func FilterByLabelID(labelID primitive.ID) filter {
	return filter{Equal, "labelid", labelID}
}

func FilterByLabel(labelOrString any) filter {
	label, err := primitive.LabelFromLabelOrString(labelOrString)
	if err != nil {
		return filter{Equal, "label", "error: " + err.Error()}
	}

	return filter{Equal, "label", label}
}

func (db *DB) matches(s adapter.Storeable, filters ...filter) bool {
	if len(filters) == 0 {
		return false
	}

	for _, f := range filters {
		if f.field == "label" {
			id, _ := db.GetLabelID(f.value)

			if s.GetLabelID() == id {
				continue
			}

			return false
		}

		if f.field == "kind" {
			if f.condition == Equal && f.value.(kind.Kind) == kind.KindOf(s) {
				continue
			} else {
				if f.value.(kind.Kind) == kind.KindOf(s) {
					return false
				}
			}
		}

		v, ok := getFieldValueByName(s, f.field)
		if !ok {
			return false
		}

		switch f.condition {
		case Equal:
			if !compareAsString(v, f.value) {
				return false
			}
		case NotEqual:
			if compareAsString(v, f.value) {
				return false
			}
		case Contains:
			if !fieldContains(v, f.value) {
				return false
			}
		case NotContains:
			if fieldContains(v, f.value) {
				return false
			}
		}
	}

	return true
}

func fieldContains(field reflect.Value, value any) bool {
	kind := field.Kind()

	switch kind {
	case reflect.String:
		return strings.Contains(field.String(), value.(string))
	case reflect.Array, reflect.Slice:
		for i := 0; i < field.Len(); i++ {
			if compareAsString(field.Index(i), value) {
				return true
			}
		}
	}

	return false
}

func getFieldValueByName(s adapter.Storeable, fieldName string) (reflect.Value, bool) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	field := val.FieldByName(fieldName)
	if field.IsValid() {
		return field, true
	}

	convertFuncs := []func(string) string{
		strings.ToLower,
		strings.ToUpper,
		strings.ToTitle,
		func(s string) string { return cases.Title(language.English).String(s) },
	}

	fieldName = strings.TrimSpace(fieldName)

	for _, convertFunc := range convertFuncs {
		name := convertFunc(fieldName)
		field = val.FieldByName(name)
		if field.IsValid() {
			return field, true
		}
	}

	return reflect.Value{}, false
}

func compareAsString(field reflect.Value, value any) bool {
	return fmt.Sprintf("%v", field.Interface()) == fmt.Sprintf("%v", value)
}
