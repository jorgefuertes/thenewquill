package validator

import (
	"fmt"
	"reflect"
	"strconv"
)

func validateMin(t reflect.StructField, v reflect.Value, params string) error {
	min, err := strconv.Atoi(params)
	if err != nil {
		return err
	}

	numErr := fmt.Errorf("%s must be at least %d", t.Name, min)
	strErr := fmt.Errorf("%s minimun length must be %d", t.Name, min)

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() < int64(min) {
			return numErr
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() < uint64(min) {
			return numErr
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() < float64(min) {
			return numErr
		}
	case reflect.String:
		if len(v.String()) < min {
			return strErr
		}
	default:
		return fmt.Errorf("unsupported type %q", v.Kind())
	}

	return nil
}

func validateMax(t reflect.StructField, v reflect.Value, params string) error {
	max, err := strconv.Atoi(params)
	if err != nil {
		return err
	}

	numErr := fmt.Errorf("%s must be %d as maximum", t.Name, max)
	strErr := fmt.Errorf("%s maximun length must be %d", t.Name, max)

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() > int64(max) {
			return numErr
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() > uint64(max) {
			return numErr
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() > float64(max) {
			return numErr
		}
	case reflect.String:
		if len(v.String()) > max {
			return strErr
		}
	default:
		return fmt.Errorf("unsupported type %q", v.Kind())
	}

	return nil
}
