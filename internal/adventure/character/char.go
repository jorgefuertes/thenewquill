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
		Description: "",
		Location:    nil,
		Created:     false,
		Human:       false,
		Vars:        vars.NewStore(),
	}
}

func (c Character) export() map[string]any {
	data := map[string]any{
		"label":       c.Label,
		"name":        c.Name.Label,
		"adjective":   c.Adjective.Label,
		"description": c.Description,
		"location":    c.Location.Label,
		"created":     c.Created,
		"human":       c.Human,
	}

	for k, v := range c.Vars.GetAll() {
		data["var:"+k] = v
	}

	return data
}
