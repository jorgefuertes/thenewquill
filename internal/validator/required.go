package validator

import (
	"fmt"
	"reflect"
)

func validateRequired(t reflect.StructField, v reflect.Value, _ string) error {
	if reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface()) {
		return fmt.Errorf("%s is required", t.Name)
	}

	return nil
}
