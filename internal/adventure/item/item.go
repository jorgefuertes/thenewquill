package item

import (
	"thenewquill/internal/adventure/character"
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/vars"
	"thenewquill/internal/adventure/words"
)

type Item struct {
	Label       string
	Noun        *words.Word
	Adjective   *words.Word
	Description string
	Weight      int
	MaxWeight   int
	IsContainer bool
	IsWearable  bool
	IsWorn      bool
	IsCreated   bool
	IsHeld      bool
	Location    *loc.Location
	CarriedBy   *character.Character
	Contents    []*Item
	Vars        vars.Store
}

// simple Item
func New(label string, noun *words.Word, adjective *words.Word) *Item {
	return &Item{
		Label:     label,
		Noun:      noun,
		Adjective: adjective,
		Weight:    0,
		MaxWeight: 100,
		Contents:  make([]*Item, 0),
		Vars:      vars.NewStore(),
	}
}

func (i Item) String() string {
	if i.Adjective != nil {
		return i.Noun.Label + " " + i.Adjective.Label
	}

	return i.Noun.Label
}

func (i *Item) Wear() {
	if i.IsWearable {
		i.IsWorn = true
		i.IsHeld = false
		i.Location = nil
	}
}

func (i *Item) Unwear() {
	i.IsWorn = false
	i.Hold()
	i.Location = nil
}

func (i *Item) Hold() {
	i.IsWorn = false
	i.IsHeld = true
}

func (i *Item) Create() {
	i.IsCreated = true
}

func (i *Item) Destroy() {
	i.IsHeld = false
	i.IsWorn = false
	i.IsCreated = false
	i.Location = nil
}

func (i *Item) SetCarriedBy(p *character.Character) {
	i.CarriedBy = p
}
