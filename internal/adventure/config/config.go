package config

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
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
	{"date", true},
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
	ID db.ID
	V  string
}

var _ db.Storeable = Param{}

func (v Param) Export() string {
	return fmt.Sprintf("%d|%d|%s\n", v.GetKind().Byte(), v.ID, v.V)
}

func (v Param) GetID() db.ID {
	return v.ID
}

func (v Param) GetKind() kind.Kind {
	return kind.Param
}

func (v Param) SetID(id db.ID) db.Storeable {
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
