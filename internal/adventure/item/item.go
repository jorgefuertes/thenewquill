package item

import (
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/player"
	"thenewquill/internal/adventure/vars"
	"thenewquill/internal/adventure/voc"
)

type Item struct {
	label       string
	noun        *voc.Word
	adjective   *voc.Word
	description string
	weight      int
	maxWeight   int
	isContainer bool
	isWearable  bool
	isWorn      bool
	isCreated   bool
	isHeld      bool
	location    *loc.Location
	carriedBy   *player.Player
	contents    []*Item
	Vars        vars.Store
}

// simple Item
func New(label string, noun *voc.Word, adjective *voc.Word) *Item {
	return &Item{
		label:     label,
		noun:      noun,
		adjective: adjective,
		weight:    0,
		maxWeight: 100,
		contents:  make([]*Item, 0),
		Vars:      vars.NewStore(),
	}
}

func (i Item) String() string {
	if i.adjective != nil {
		return i.noun.Label + " " + i.adjective.Label
	}

	return i.noun.Label
}

func (i Item) Label() string {
	return i.label
}

func (i *Item) SetContainer() {
	i.isContainer = true
}

func (i *Item) SetWeight(w int) {
	i.weight = w
}

func (i *Item) SetMaxWeight(w int) {
	i.maxWeight = w
}

func (i *Item) SetDescription(d string) {
	i.description = d
}

func (i *Item) Description() string {
	return i.description
}

func (i *Item) SetWearable() {
	i.isWearable = true
}

func (i *Item) IsWearable() bool {
	return i.isWearable
}

func (i Item) IsWorn() bool {
	return i.isWorn
}

func (i *Item) Wear() {
	if i.isWearable {
		i.isWorn = true
		i.isHeld = false
		i.location = nil
	}
}

func (i *Item) Unwear() {
	i.isWorn = false
	i.Hold()
	i.location = nil
}

func (i Item) IsHeld() bool {
	return i.isHeld
}

func (i *Item) Hold() {
	i.isWorn = false
	i.isHeld = true
}

func (i Item) IsCreated() bool {
	return i.isCreated
}

func (i *Item) Create() {
	i.isCreated = true
}

func (i *Item) Destroy() {
	i.isHeld = false
	i.isWorn = false
	i.isCreated = false
	i.location = nil
}

func (i *Item) SetCarriedBy(p *player.Player) {
	i.carriedBy = p
}

func (i *Item) CarriedBy() *player.Player {
	return i.carriedBy
}

func (i *Item) SetLocation(l *loc.Location) {
	i.location = l
}

func (i *Item) Location() *loc.Location {
	return i.location
}
