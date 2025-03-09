package character

import (
	"fmt"

	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/vars"
	"thenewquill/internal/adventure/words"
	"thenewquill/internal/compiler/section"
	"thenewquill/internal/util"
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

func (c Character) export() []string {
	locationLabel := ""
	if c.Location != nil {
		locationLabel = c.Location.Label
	}

	data := []string{
		util.ValueToString(section.Chars.ToInt()),
		c.Label,
		c.Name.Label,
		c.Adjective.Label,
		c.Description,
		locationLabel,
		util.ValueToString(c.Created),
		util.ValueToString(c.Human),
	}

	for k, v := range c.Vars.GetAll() {
		data = append(data, fmt.Sprintf("var:%s=%s", k, util.ValueToString(v)))
	}

	return data
}
