package validator

import (
	"fmt"
	"reflect"
	"regexp"
)

type validator struct {
	tag string
	rg  *regexp.Regexp
	fn  func(t reflect.StructField, v reflect.Value, params string) error
}

var validators = []validator{
	{tag: "required", rg: nil, fn: validateRequired},
	{tag: "numeric", rg: nil, fn: validateNumeric},
	{tag: "min", rg: regexp.MustCompile(`^min=(\d+)$`), fn: validateMin},
	{tag: "max", rg: regexp.MustCompile(`^max=(\d+)$`), fn: validateMax},
	{tag: "matches", rg: regexp.MustCompile(`^matches\((.+)\)$`), fn: validateRegexp},
	{tag: "in", rg: regexp.MustCompile(`^in=(.+)$`), fn: validateIn},
	{tag: "len", rg: regexp.MustCompile(`^len\((\d+(?:\|\d+)?)\)$`), fn: validateLen},
	{tag: "count", rg: regexp.MustCompile(`^count\((\d+(?:\|\d+)?)\)$`), fn: validateCount},
}

func Validate(input any) error {
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		valid, ok := t.Field(i).Tag.Lookup("valid")
		if !ok {
			continue
		}

		for _, tag := range parseValidTags(valid) {
			validator, params, err := getValidatorAndParams(tag)
			if err != nil {
				return err
			}

			if err := validator.fn(t.Field(i), v.Field(i), params); err != nil {
				return fmt.Errorf("validating %s: %s", t.Name(), err)
			}
		}
	}

	return nil
}
