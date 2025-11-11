package config

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

type allowed struct {
	labelName string
	required  bool
}

var allowedFields = []allowed{
	{"title", true},
	{"author", true},
	{"description", true},
	{"version", true},
	{"date", false},
	{"language", true},
}

func AllowedFieldNames() []string {
	fields := make([]string, 0)

	for _, allowed := range allowedFields {
		fields = append(fields, allowed.labelName)
	}

	return fields
}

type Param struct {
	ID id.ID
	V  string
}

var _ adapter.Storeable = Param{}

func (v Param) GetID() id.ID {
	return v.ID
}

func (v Param) GetKind() kind.Kind {
	return kind.Param
}

func (v Param) SetID(id id.ID) adapter.Storeable {
	v.ID = id

	return v
}

func isKeyAllowed(key string) bool {
	for _, allowed := range allowedFields {
		if key == allowed.labelName {
			return true
		}
	}

	return false
}
