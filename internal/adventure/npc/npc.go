package npc

import (
	"thenewquill/internal/adventure/item"
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/voc"
)

type NPC struct {
	label       string
	Name        *voc.Word
	Description string
	Location    *loc.Location
	Inventory   []*item.Item
	Created     bool
}

func New(label string, name *voc.Word) *NPC {
	return &NPC{
		label:       label,
		Name:        name,
		Description: "",
		Location:    nil,
		Inventory:   make([]*item.Item, 0),
		Created:     false,
	}
}
