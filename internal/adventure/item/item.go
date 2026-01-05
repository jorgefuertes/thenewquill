package item

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database/adapter"
)

type Item struct {
	ID          uint32
	LabelID     uint32 `valid:"required"`
	NounID      uint32 `valid:"required"`
	AdjectiveID uint32
	Description string `valid:"required"`
	Weight      int
	MaxWeight   int
	Container   bool
	Wearable    bool
	Created     bool
	At          uint32
	Worn        bool
}

var _ adapter.Storeable = &Item{}

const defaultMaxWeight = 100

func New() *Item {
	return &Item{MaxWeight: defaultMaxWeight}
}

func (i *Item) GetKind() kind.Kind {
	return kind.Item
}

func (i *Item) SetID(id uint32) {
	i.ID = id
}

func (i *Item) GetID() uint32 {
	return i.ID
}

func (i *Item) SetLabelID(id uint32) {
	i.LabelID = id
}

func (i *Item) GetLabelID() uint32 {
	return i.LabelID
}
