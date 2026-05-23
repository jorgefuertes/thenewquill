package process

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database/adapter"
)

type TableKind uint8

/*
 * 0: init     Cuando arranca la aventura, o se reinicia
 * 1: location Cuando cambia de ubicación
 * 2: turn     Tras el input, antes de `item`
 * 3: item     Tras el input, sólo si la SL contiene un item
 * 4: npc      Tras el input, sólo si la SL contiene un NPC
 * 5: response Tras el input, en último lugar
 * 6: cron     Por tiempo o turnos, procesos independientes
 */

const (
	Init TableKind = iota
	Location
	Turn
	Item
	NPC
	Response
	Cron
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
