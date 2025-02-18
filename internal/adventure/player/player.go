package player

import (
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/voc"
)

type Player struct {
	label       string
	Name        *voc.Word
	Description string
	Location    *loc.Location
	Created     bool
	Human bool
}

func New(label string, name *voc.Word) *Player {
	return &Player{
		label:       label,
		Name:        name,
		Description: "",
		Location:    nil,
		Created:     false,
		Human: false,
	}
}

func NewHuman(label string) *Player {
	return &Player{
		label:       label,
		Name:        nil,
		Description: "",
		Location:    nil,
		Created:     false,
		Human: true,
	}
}
