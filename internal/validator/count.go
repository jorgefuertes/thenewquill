package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func validateCount(t reflect.StructField, v reflect.Value, params string) error {
	var min int
	var max int
	var err error

	if strings.Contains(params, "|") {
		parts := strings.Split(params, "|")
		min, err = strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("%s: invalid min count %q", t.Name, params)
		}

		max, err = strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("%s: invalid max count` %q", t.Name, params)
		}
	} else {
		min, err = strconv.Atoi(params)
		if err != nil {
			return fmt.Errorf("%s: invalid min count %q", t.Name, params)
		}
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		if min > 0 && v.Len() < min {
			return fmt.Errorf("%s must have at least %d elements", t.Name, min)
		}

		if max > 0 && v.Len() > max {
			return fmt.Errorf("%s must have less or equal than %d elements", t.Name, max)
		}
	default:
		return fmt.Errorf("unsupported type %q", v.Kind())
	}

	return nil
}
