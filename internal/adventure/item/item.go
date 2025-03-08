package item

import (
	"errors"
	"fmt"
	"strings"

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

func (i *Item) Validate() error {
	if i.Noun == nil {
		return ErrNounCannotBeNil
	}

	if i.Noun.Is("_") {
		return ErrNounCannotBeUnderscore
	}

	if i.Adjective == nil {
		return ErrAdjectiveCannotBeNil
	}

	if i.IsHeld && i.Location != nil {
		return ErrItemCannotBeHeldAndHaveLocation
	}

	if i.CarriedBy != nil && i.Location != nil {
		return ErrItemCannotBeHeldAndHaveLocation
	}

	if i.IsHeld && i.IsWorn {
		return ErrItemCannotBeHeldAndWorn
	}

	for _, content := range i.Contents {
		if content.Location != nil {
			return errors.Join(ErrItemCannotBeContainedInAndHaveLocation,
				fmt.Errorf("item %s is at %s and contained in %s", content.Label, content.Location.Label, i.Label))
		}
	}

	if i.IsContainer && i.WeightTotal() > i.MaxWeight {
		return ErrContainerCantCarrySoMuch
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

func (i Item) export() map[string]any {
	data := map[string]any{
		"label":       i.Label,
		"noun":        i.Noun.Label,
		"adjective":   i.Adjective.Label,
		"description": i.Description,
		"weight":      i.Weight,
		"maxWeight":   i.MaxWeight,
		"isContainer": i.IsContainer,
		"isWearable":  i.IsWearable,
		"isWorn":      i.IsWorn,
		"isCreated":   i.IsCreated,
		"isHeld":      i.IsHeld,
		"location":    i.Location.Label,
		"carriedBy":   i.CarriedBy.Label,
	}

	contentLabels := make([]string, 0)
	for _, c := range i.Contents {
		contentLabels = append(contentLabels, c.Label)
	}
	data["contents"] = strings.Join(contentLabels, ",")

	for k, v := range i.Vars.GetAll() {
		data["var:"+k] = v
	}

	return data
}
