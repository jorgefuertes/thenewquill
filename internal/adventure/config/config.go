package config

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
)

type allowed struct {
	label    primitive.Label
	required bool
}

var allowedFieldLabels = []allowed{
	{"title", true},
	{"author", true},
	{"description", true},
	{"version", true},
	{"date", false},
	{"language", true},
}

func AllowedFieldLabels() []primitive.Label {
	fields := make([]primitive.Label, 0)

	for _, allowed := range allowedFieldLabels {
		fields = append(fields, allowed.label)
	}

	return fields
}

type Param struct {
	ID      primitive.ID
	LabelID primitive.ID
	V       string
}

func New(id, labelID primitive.ID, v string) *Param {
	return &Param{ID: id, V: v}
}

var _ adapter.Storeable = &Param{}

func (v Param) GetID() primitive.ID {
	return v.ID
}

func (v *Param) SetID(id primitive.ID) {
	v.ID = id
}

func (v Param) GetLabelID() primitive.ID {
	return v.LabelID
}

func (v *Param) SetLabelID(id primitive.ID) {
	v.LabelID = id
}

func isLabelAllowed(label primitive.Label) bool {
	for _, allowed := range allowedFieldLabels {
		if label == allowed.label {
			return true
		}
	}

	return false
}
