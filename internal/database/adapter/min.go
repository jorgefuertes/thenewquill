package adapter

import "github.com/jorgefuertes/thenewquill/internal/adventure/kind"

type MinStoreable struct {
	ID      uint32
	LabelID uint32
	Kind    kind.Kind
}

var _ Storeable = &MinStoreable{}

func (m *MinStoreable) SetID(id uint32) {
	m.ID = id
}

func (m *MinStoreable) GetID() uint32 {
	return m.ID
}

func (m *MinStoreable) SetLabelID(labelID uint32) {
	m.LabelID = labelID
}

func (m *MinStoreable) GetLabelID() uint32 {
	return m.LabelID
}

func (m *MinStoreable) GetKind() kind.Kind {
	return m.Kind
}

func (m *MinStoreable) SetKind(kind kind.Kind) {
	m.Kind = kind
}
