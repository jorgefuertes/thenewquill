package character

import (
	"thenewquill/internal/adventure/loc"
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
}

func New(label string, name *words.Word, adjective *words.Word) *Character {
	return &Character{
		Label:       label,
		Name:        name,
		Description: "",
		Location:    nil,
		Created:     false,
		Human:       false,
	}
}

func NewHuman(label string) *Character {
	return &Character{
		Label:       label,
		Name:        nil,
		Description: "",
		Location:    nil,
		Created:     false,
		Human:       true,
	}
}
