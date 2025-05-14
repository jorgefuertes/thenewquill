package item

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

type Item struct {
	ID          db.ID
	NounID      db.ID
	AdjectiveID db.ID
	Description string
	Weight      int
	MaxWeight   int
	IsContainer bool
	Contains    []db.ID
	IsWearable  bool
	WornBy      db.ID
	IsCreated   bool
	LocationID  db.ID
}

var _ db.Storeable = Item{}

func (i Item) GetID() db.ID {
	return i.ID
}

func (i Item) GetKind() (db.Kind, db.SubKind) {
	return db.Items, db.NoSubKind
}

// simple Item
func New(id db.ID, nounID db.ID, adjectiveID db.ID) Item {
	return Item{
		ID:          id,
		NounID:      nounID,
		AdjectiveID: adjectiveID,
		Weight:      0,
		MaxWeight:   100,
	}
}
