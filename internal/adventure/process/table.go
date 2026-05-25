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

func AllowedTables() []TableKind {
	return []TableKind{
		Init,
		Location,
		Turn,
		Item,
		NPC,
		Response,
		Cron,
	}
}

func (k TableKind) String() string {
	switch k {
	case Init:
		return "init"
	case Location:
		return "location"
	case Turn:
		return "turn"
	case Item:
		return "item"
	case NPC:
		return "npc"
	case Response:
		return "response"
	case Cron:
		return "cron"
	default:
		return ""
	}
}

type Table struct {
	ID      uint32
	LabelID uint32
	Kind    TableKind
	Procs   []Process
}

var _ adapter.Storeable = &Table{}

func NewTable(kind TableKind) *Table {
	m := &Table{Kind: kind}

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
