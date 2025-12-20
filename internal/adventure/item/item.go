package item

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
)

type Item struct {
	ID          uint32
	LabelID     uint32
	NounID      uint32
	AdjectiveID uint32
	Description string
	Weight      int
	MaxWeight   int
	Container   bool
	Wearable    bool
	Created     bool
	At          uint32
	Worn        bool
}

var _ adapter.Storeable = &Item{}

func New(id, nounID, adjectiveID uint32) *Item {
	return &Item{
		ID:          id,
		NounID:      nounID,
		AdjectiveID: adjectiveID,
		Weight:      0,
		MaxWeight:   100,
	}
}

func (i *Item) SetID(id uint32) {
	i.ID = id
}

func (i Item) GetID() uint32 {
	return i.ID
}

func (i *Item) SetLabelID(id uint32) {
	i.LabelID = id
}

func (i Item) GetLabelID() uint32 {
	return i.LabelID
}
