package bin

import (
	"thenewquill/internal/adventure"
	"thenewquill/internal/adventure/config"
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/words"
)

const DATABASE_VERSION = "1"

type adventureDTO struct {
	Config    config.Config
	Labels    labels
	Vars      map[string]any
	Msgs      []msgDTO
	Words     []wordDTO
	Locations []locationDTO
	Items     []itemDTO
	Chars     []charDTO
}

type msgDTO struct {
	ID      int
	Text    string
	Plurals [3]string
}

type wordDTO struct {
	ID       int
	Type     words.WordType
	Synonyms []string
}

type locationDTO struct {
	ID          int
	Title       string
	Description string
	Conns       map[int]int
}

type itemDTO struct {
	ID          int
	NounID      int
	AdjectiveID int
	Weight      int
	MaxWeight   int
	IsContainer bool
	IsWearable  bool
	IsWorn      bool
	IsCreated   bool
	IsHeld      bool
	LocationID  int
	CarriedByID int
}

type charDTO struct {
	ID          int
	NameID      int
	AdjectiveID int
	Description string
	LocationID  int
	Created     bool
	Human       bool
}

func (d *adventureDTO) GetWordByLabelID(id int, a *adventure.Adventure) *words.Word {
	l, ok := d.Labels.Get(id)
	if !ok {
		return nil
	}

	labelName, t := splitWordName(l.Name)
	if t == words.Unknown {
		return nil
	}

	return a.Words.Get(t, labelName)
}

func (d *adventureDTO) GetLocationByLabelID(id int, a *adventure.Adventure) *loc.Location {
	l, ok := d.Labels.Get(id)
	if !ok {
		return nil
	}

	return a.Locations.Get(l.Name)
}
