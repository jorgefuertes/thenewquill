package db

import (
	"reflect"
	"strings"
)

type filter struct {
	field string
	value any
}

func Filter(field string, value any) filter {
	return filter{field, value}
}

func FilterByKind(kind Kind) filter {
	return filter{"kind", kind}
}

func FilterByID(id ID) filter {
	return filter{"id", id}
}

func matches(s Storeable, filters ...filter) bool {
	if len(filters) == 0 {
		return false
	}

	for _, f := range filters {
		switch strings.ToLower(f.field) {
		case "id":
			if s.GetID() != f.value.(ID) {
				return false
			}
		case "kind":
			if s.GetKind() != f.value.(Kind) {
				return false
			}
		default:
			field := reflect.ValueOf(s).FieldByName(f.field)

			if !field.IsValid() || field.Interface() != f.value {
				return false
			}
		}
	}

	return true
}
