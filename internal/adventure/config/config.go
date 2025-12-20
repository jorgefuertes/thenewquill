package config

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/database/primitive"
)

type Param struct {
	ID      uint32
	LabelID uint32
	V       string
}

func New(id, labelID uint32, v string) *Param {
	return &Param{ID: id, V: v}
}

var _ adapter.Storeable = &Param{}

func (v Param) GetID() uint32 {
	return v.ID
}

func (v *Param) SetID(id uint32) {
	v.ID = id
}

func (v Param) GetLabelID() uint32 {
	return v.LabelID
}

func (v *Param) SetLabelID(id uint32) {
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
