package config

import "github.com/jorgefuertes/thenewquill/internal/adventure/db"

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
	{"lang", true},
}

func AllowedFieldNames() []string {
	fields := make([]string, 0)

	for _, allowed := range allowedFields {
		fields = append(fields, allowed.labelName)
	}

	return fields
}

type Value struct {
	ID db.ID
	V  string
}

var _ db.Storeable = Value{}

func (v Value) GetID() db.ID {
	return v.ID
}

func (v Value) GetKind() db.Kind {
	return db.Config
}

func (v Value) SetID(id db.ID) db.Storeable {
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
