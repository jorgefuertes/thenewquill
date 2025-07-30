package validator

import (
	"fmt"
	"reflect"
	"regexp"
)

func validateRegexp(t reflect.StructField, v reflect.Value, params string) error {
	rg, err := regexp.Compile(params)
	if err != nil {
		return err
	}

	if v.Kind() != reflect.String {
		return fmt.Errorf("regexp tag %q is only supported on string fields", t.Name)
	}

	if v.String() == "" {
		return nil
	}

	if rg.MatchString(v.String()) {
		return nil
	}

	return fmt.Errorf("%s does not match regexp `%v`", t.Name, params)
}
