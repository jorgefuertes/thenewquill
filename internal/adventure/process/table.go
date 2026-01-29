package process

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database/adapter"
)

type TableKind uint8

const (
	Init TableKind = iota
	Location
	Response
	Cron
	NPC
)

type Table struct {
	ID         uint32
	LabelID    uint32
	Kind       TableKind
	ProcessIDs []uint32
}

var _ adapter.Storeable = &Table{}

func NewTable() *Table {
	m := &Table{}

	return m
}

func (m *Table) GetKind() kind.Kind {
	return kind.Table
}

func (m *Table) SetID(id uint32) {
	m.ID = id
}

func (m *Table) GetID() uint32 {
	return m.ID
}

func (m *Table) SetLabelID(id uint32) {
	m.LabelID = id
}

func (m *Table) GetLabelID() uint32 {
	return m.LabelID
}
