package db

import (
	"fmt"
	"reflect"
	"strings"

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

func FilterByID(id ID) filter {
	return filter{Equal, "id", id}
}

func matches(s Storeable, filters ...filter) bool {
	if len(filters) == 0 {
		return false
	}

	for _, f := range filters {
		if f.field == "kind" {
			if f.condition == Equal && f.value.(kind.Kind) == kind.KindOf(s) {
				return true
			} else {
				return f.value.(kind.Kind) != kind.KindOf(s)
			}
		}

		field, ok := getFieldByName(s, f.field)
		if !ok {
			return false
		}

		switch f.condition {
		case Equal:
			if !compareAsString(field, f.value) {
				return false
			}
		case NotEqual:
			if compareAsString(field, f.value) {
				return false
			}
		case Contains:
			if !fieldContains(field, f.value) {
				return false
			}
		case NotContains:
			if fieldContains(field, f.value) {
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

func getFieldByName(s Storeable, fieldName string) (reflect.Value, bool) {
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
