package character

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
)

type Character struct {
	ID          id.ID
	NounID      id.ID
	AdjectiveID id.ID
	Description string
	LocationID  id.ID
	Created     bool
	Human       bool
}

var _ adapter.Storeable = &Character{}

func New(nounID, adjectiveID id.ID) Character {
	return Character{
		ID:          id.Undefined,
		NounID:      nounID,
		AdjectiveID: adjectiveID,
		Description: "",
		LocationID:  id.Undefined,
		Created:     false,
		Human:       false,
	}
}

func (c Character) SetID(i id.ID) adapter.Storeable {
	c.ID = i

	return c
}

func (c Character) GetID() id.ID {
	return c.ID
}
