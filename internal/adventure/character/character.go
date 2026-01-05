package character

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database/adapter"
)

type Character struct {
	ID          uint32
	LabelID     uint32 `valid:"required"`
	NounID      uint32 `valid:"required"`
	AdjectiveID uint32 `valid:"required"`
	Description string `valid:"required"`
	LocationID  uint32 `valid:"required"`
	Created     bool
	Human       bool
}

var _ adapter.Storeable = &Character{}

func New() *Character {
	return &Character{}
}

func (c *Character) GetKind() kind.Kind {
	return kind.Character
}

func (c *Character) SetID(id uint32) {
	c.ID = id
}

func (c *Character) GetID() uint32 {
	return c.ID
}

func (c *Character) SetLabelID(id uint32) {
	c.LabelID = id
}

func (c *Character) GetLabelID() uint32 {
	return c.LabelID
}
