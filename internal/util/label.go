package util

import "reflect"

type Labeled interface {
	GetLabel() string
}

func ToLabel(l Labeled) string {
	if l == nil || reflect.ValueOf(l).IsNil() {
		return ""
	}

	return l.GetLabel()
}
