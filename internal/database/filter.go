package database

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fxamacker/cbor/v2"
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
	Matches
	NotMatches
)

const (
	idFieldName      = "ID"
	labelIDFieldName = "LabelID"
	labelFieldName   = "Label"
	kindFieldName    = "Kind"
)

type Filter struct {
	condition condition
	field     string
	value     any
}

func NewFilter(field string, condition condition, value any) Filter {
	return Filter{condition, field, value}
}

func FilterByKind(k kind.Kind) Filter {
	return Filter{Equal, kindFieldName, k}
}

func FilterByID(id uint32) Filter {
	return Filter{Equal, idFieldName, id}
}

func FilterByLabelID(labelID uint32) Filter {
	return Filter{Equal, labelIDFieldName, labelID}
}

func FilterByLabel(label string) Filter {
	return Filter{Equal, labelFieldName, label}
}

func (db *DB) matchesAllFilters(r Record, filters ...Filter) bool {
	if len(filters) == 0 {
		return false
	}

	recordMap := map[string]any{}
	if err := cbor.Unmarshal(r.Data, &recordMap); err != nil {
		return false
	}

	for _, f := range filters {
		if strings.EqualFold(f.field, kindFieldName) {
			if f.condition == Equal && f.value.(kind.Kind) == r.Kind {
				continue
			} else {
				if f.value.(kind.Kind) == kind.KindOf(r) {
					return false
				}
			}
		}

		if strings.EqualFold(f.field, labelFieldName) {
			labelID, err := db.getLabelID(f.value.(string))
			if err == nil {
				if f.condition == Equal && labelID == r.LabelID {
					continue
				} else {
					if labelID == r.LabelID {
						return false
					}
				}
			}
		}

		if !checkCondition(recordMap, f) {
			return false
		}
	}

	return true
}

func getField(recordMap map[string]any, fieldName string) (any, bool) {
	v, ok := recordMap[fieldName]
	if ok {
		return v, true
	}

	titleCaser := cases.Title(language.Und)
	conversors := []func(string) string{
		strings.ToLower,
		strings.ToUpper,
		titleCaser.String,
	}

	for _, c := range conversors {
		v, ok = recordMap[c(fieldName)]
		if ok {
			return v, true
		}
	}

	return nil, false
}

func checkCondition(recordMap map[string]any, f Filter) bool {
	v, ok := getField(recordMap, f.field)
	if !ok {
		return false
	}

	switch f.condition {
	case Equal:
		if !compare(v, f.value) {
			return false
		}
	case NotEqual:
		if compare(v, f.value) {
			return false
		}
	case Contains:
		if !contains(v, f.value) {
			return false
		}
	case NotContains:
		if contains(v, f.value) {
			return false
		}
	case Matches:
		rg := regexp.MustCompile(f.value.(string))
		if !rg.MatchString(fmt.Sprintf("%v", v)) {
			return false
		}
	case NotMatches:
		rg := regexp.MustCompile(f.value.(string))
		if rg.MatchString(fmt.Sprintf("%v", v)) {
			return false
		}
	}

	return true
}

func contains(haystack, needle any) bool {
	switch hs := haystack.(type) {
	case string:
		return strings.Contains(hs, needle.(string))
	case []any:
		for _, v := range hs {
			if compare(v, needle) {
				return true
			}
		}
	}

	return false
}

func compare(v1, v2 any) bool {
	switch v1.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		return fmt.Sprintf("%d", v1) == fmt.Sprintf("%d", v2)
	case float32, float64:
		return fmt.Sprintf("%.4f", v1) == fmt.Sprintf("%.4f", v2)
	case bool:
		return fmt.Sprintf("%t", v1) == fmt.Sprintf("%t", v2)
	case string:
		return fmt.Sprintf("%s", v1) == fmt.Sprintf("%s", v2)
	default:
		return fmt.Sprintf("%v", v1) == fmt.Sprintf("%v", v2)
	}
}
