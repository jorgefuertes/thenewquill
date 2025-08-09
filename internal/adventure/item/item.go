package item

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

type Item struct {
	ID          db.ID
	NounID      db.ID
	AdjectiveID db.ID
	Description string
	Weight      int
	MaxWeight   int
	Container   bool
	Wearable    bool
	Created     bool
	At          db.ID
	Worn        bool
}

var _ db.Storeable = Item{}

func New(nounID db.ID, adjectiveID db.ID) Item {
	return Item{
		ID:          db.UndefinedLabel.ID,
		NounID:      nounID,
		AdjectiveID: adjectiveID,
		Weight:      0,
		MaxWeight:   100,
	}
}

func (i Item) SetID(id db.ID) db.Storeable {
	i.ID = id

	return i
}

func (i Item) GetID() db.ID {
	return i.ID
}

func (i Item) GetKind() kind.Kind {
	return kind.Item
}
