package validator

import (
	"fmt"
	"reflect"
	"regexp"
)

func validateNumeric(t reflect.StructField, v reflect.Value, _ string) error {
	switch v.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64:
		return nil
	case reflect.String:
		if v.String() == "" {
			return nil
		}

		if regexp.MustCompile(`^\d+$`).MatchString(v.String()) {
			return nil
		}

		return fmt.Errorf("%s must be numeric", t.Name)
	default:
		return fmt.Errorf("unsupported type %q", v.Kind())
	}
}
