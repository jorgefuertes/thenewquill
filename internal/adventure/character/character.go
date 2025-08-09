package character

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

type Character struct {
	ID          db.ID
	NounID      db.ID
	AdjectiveID db.ID
	Description string
	LocationID  db.ID
	Created     bool
	Human       bool
}

var _ db.Storeable = &Character{}

func New(nounID db.ID, adjectiveID db.ID) Character {
	return Character{
		ID:          db.UndefinedLabel.ID,
		NounID:      nounID,
		AdjectiveID: adjectiveID,
		Description: "",
		LocationID:  db.UndefinedLabel.ID,
		Created:     false,
		Human:       false,
	}
}

func (c Character) SetID(id db.ID) db.Storeable {
	c.ID = id

	return c
}

func (c Character) GetID() db.ID {
	return c.ID
}
