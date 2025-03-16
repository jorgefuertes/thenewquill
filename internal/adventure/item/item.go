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
	IsCreated   bool
	Location    *loc.Location
	Inside      *Item
	CarriedBy   *character.Character
	WornBy      *character.Character
	Vars        vars.Store
}

func (i *Item) validate() error {
	if i.Noun == nil {
		return ErrNounCannotBeNil
	}

	if i.Noun.Is("_") {
		return ErrNounCannotBeUnderscore
	}

	if i.Adjective == nil {
		return ErrAdjectiveCannotBeNil
	}

	if i.Weight > i.MaxWeight {
		return ErrWeightShouldBeLessOrEqualThanMaxWeight
	}

	if i.Weight < 0 || i.MaxWeight < 0 {
		return ErrWeightCannotBeNegative
	}

	return nil
}

// simple Item
func New(label string, noun *words.Word, adjective *words.Word) *Item {
	return &Item{
		Label:     label,
		Noun:      noun,
		Adjective: adjective,
		Weight:    0,
		MaxWeight: 100,
		Vars:      vars.NewStore(),
	}
}

func (i *Item) GetLabel() string {
	return i.Label
}

func (i *Item) Wear() {
	if i.IsWearable {
		i.WornBy = i.CarriedBy
		i.CarriedBy = nil
		i.Location = nil
		i.Inside = nil
	}
}

func (i *Item) Unwear() {
	i.CarriedBy = i.WornBy
	i.WornBy = nil
	i.Location = nil
}

func (i *Item) Drop(l *loc.Location) {
	i.WornBy = nil
	i.Inside = nil
	i.CarriedBy = nil
	i.Location = l
}

func (i *Item) Give(c *character.Character) {
	i.WornBy = nil
	i.Location = nil
	i.Inside = nil
	i.CarriedBy = c
}

func (i *Item) Create() {
	i.WornBy = nil
	i.Inside = nil
	i.CarriedBy = nil
	i.IsCreated = true
}

func (i *Item) Destroy() {
	i.WornBy = nil
	i.Inside = nil
	i.CarriedBy = nil
	i.IsCreated = false
}
