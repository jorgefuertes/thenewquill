package config

import "github.com/jorgefuertes/thenewquill/internal/database/adapter"

type Param struct {
	ID      uint32
	LabelID uint32 `valid:"required"`
	V       string `valid:"required"`
}

var _ adapter.Storeable = &Param{}

func New(id, labelID uint32, v string) *Param {
	return &Param{ID: id, V: v}
}

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
