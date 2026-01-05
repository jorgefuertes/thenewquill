package config

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database/adapter"
)

type Param struct {
	ID      uint32
	LabelID uint32 `valid:"required"`
	V       string `valid:"required"`
}

var _ adapter.Storeable = &Param{}

func New(id, labelID uint32, v string) *Param {
	return &Param{ID: id, V: v}
}

func (p Param) GetKind() kind.Kind {
	return kind.Param
}

func (p Param) GetID() uint32 {
	return p.ID
}

func (p *Param) SetID(id uint32) {
	p.ID = id
}

func (p Param) GetLabelID() uint32 {
	return p.LabelID
}

func (p *Param) SetLabelID(id uint32) {
	p.LabelID = id
}
