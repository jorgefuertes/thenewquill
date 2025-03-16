package character

import (
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/vars"
	"thenewquill/internal/adventure/words"
)

type Character struct {
	Label       string
	Name        *words.Word
	Adjective   *words.Word
	Description string
	Location    *loc.Location
	Created     bool
	Human       bool
	Vars        vars.Store
}

func New(label string, name *words.Word, adjective *words.Word) *Character {
	return &Character{
		Label:       label,
		Name:        name,
		Adjective:   adjective,
		Description: "",
		Location:    nil,
		Created:     false,
		Human:       false,
		Vars:        vars.NewStore(),
	}
}

func (c Character) GetLabel() string {
	return c.Label
}
