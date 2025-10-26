package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func validateLen(t reflect.StructField, v reflect.Value, params string) error {
	var min int
	var max int
	var err error

	if strings.Contains(params, "|") {
		parts := strings.Split(params, "|")
		min, err = strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("%s: invalid min length %q", t.Name, params)
		}

		max, err = strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("%s: invalid max length %q", t.Name, params)
		}
	} else {
		min, err = strconv.Atoi(params)
		if err != nil {
			return fmt.Errorf("%s: invalid min length %q", t.Name, params)
		}
	}

	switch v.Kind() {
	case reflect.String:
		return checkStrLen(v.String(), min, max)
	case reflect.Slice, reflect.Array:
		a, ok := v.Interface().([]string)
		if !ok {
			return fmt.Errorf("%s: invalid type %q", t.Name, v.Kind())
		}

		for i, s := range a {
			err := checkStrLen(s, min, max)
			if err != nil {
				return fmt.Errorf("%s[%d]: %s", t.Name, i, err)
			}
		}
	case reflect.Map:
		m, ok := v.Interface().(map[string]string)
		if !ok {
			return fmt.Errorf("%s: invalid type %q", t.Name, v.Kind())
		}

		for k, s := range m {
			err := checkStrLen(s, min, max)
			if err != nil {
				return fmt.Errorf("%s[%s]: %s", t.Name, k, err)
			}
		}
	default:
		return fmt.Errorf("unsupported type %q", v.Kind())
	}

	return nil
}

func checkStrLen(s string, min, max int) error {
	if min > 0 && len(s) < min {
		return fmt.Errorf("%s must have at least %d characters", s, min)
	}

	if max > 0 && len(s) > max {
		return fmt.Errorf("%s must have less or equal than %d characters", s, max)
	}

	return nil
}
