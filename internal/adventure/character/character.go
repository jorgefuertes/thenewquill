package character

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
)

type Character struct {
	ID          primitive.ID
	LabelID     primitive.ID
	NounID      primitive.ID
	AdjectiveID primitive.ID
	Description string
	LocationID  primitive.ID
	Created     bool
	Human       bool
}

var _ adapter.Storeable = &Character{}

func New(id, labelID, nounID, adjectiveID primitive.ID) *Character {
	return &Character{
		ID:          id,
		LabelID:     labelID,
		NounID:      nounID,
		AdjectiveID: adjectiveID,
		Description: "",
		LocationID:  primitive.UndefinedID,
		Created:     false,
		Human:       false,
	}
}

func (c Character) GetID() primitive.ID {
	return c.ID
}

func (c *Character) SetID(id primitive.ID) {
	c.ID = id
}

func (c Character) GetLabelID() primitive.ID {
	return c.LabelID
}

func (c *Character) SetLabelID(id primitive.ID) {
	c.LabelID = id
}
