package item

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

type Item struct {
	ID          id.ID
	NounID      id.ID
	AdjectiveID id.ID
	Description string
	Weight      int
	MaxWeight   int
	Container   bool
	Wearable    bool
	Created     bool
	At          id.ID
	Worn        bool
}

var _ adapter.Storeable = Item{}

func New(nounID id.ID, adjectiveID id.ID) Item {
	return Item{
		ID:          id.Undefined,
		NounID:      nounID,
		AdjectiveID: adjectiveID,
		Weight:      0,
		MaxWeight:   100,
	}
}

func (i Item) SetID(id id.ID) adapter.Storeable {
	i.ID = id

	return i
}

func (i Item) GetID() id.ID {
	return i.ID
}

func (i Item) GetKind() kind.Kind {
	return kind.Item
}
