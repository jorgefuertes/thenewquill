package process

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database/adapter"
)

type Process struct {
	ID      uint32
	LabelID uint32
	VerbID  uint32
	NounID  uint32
}

var _ adapter.Storeable = &Process{}

func NewProcess() *Process {
	m := &Process{}

	return m
}

func (m *Process) GetKind() kind.Kind {
	return kind.Process
}

func (m *Process) SetID(id uint32) {
	m.ID = id
}

func (m *Process) GetID() uint32 {
	return m.ID
}

func (m *Process) SetLabelID(id uint32) {
	m.LabelID = id
}

func (m *Process) GetLabelID() uint32 {
	return m.LabelID
}
