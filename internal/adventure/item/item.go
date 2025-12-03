package item

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
)

type Item struct {
	ID          primitive.ID
	LabelID     primitive.ID
	NounID      primitive.ID
	AdjectiveID primitive.ID
	Description string
	Weight      int
	MaxWeight   int
	Container   bool
	Wearable    bool
	Created     bool
	At          primitive.ID
	Worn        bool
}

var _ adapter.Storeable = &Item{}

func New(id, nounID, adjectiveID primitive.ID) *Item {
	return &Item{
		ID:          id,
		NounID:      nounID,
		AdjectiveID: adjectiveID,
		Weight:      0,
		MaxWeight:   100,
	}
}

func (i *Item) SetID(id primitive.ID) {
	i.ID = id
}

func (i Item) GetID() primitive.ID {
	return i.ID
}

func (i *Item) SetLabelID(id primitive.ID) {
	i.LabelID = id
}

func (i Item) GetLabelID() primitive.ID {
	return i.LabelID
}
