package adventure

import (
	"thenewquill/internal/adventure/config"
	"thenewquill/internal/adventure/item"
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/msg"
	"thenewquill/internal/adventure/vars"
	"thenewquill/internal/adventure/voc"
)

type Adventure struct {
	Config     config.Config
	Vars       vars.Store
	Vocabulary voc.Vocabulary
	Messages   msg.Store
	Locations  loc.Store
	Items      item.Store
}

func New() *Adventure {
	return &Adventure{
		Config:     config.New(),
		Vars:       vars.NewStore(),
		Vocabulary: voc.NewStore(),
		Messages:   msg.NewStore(),
		Locations:  loc.NewStore(),
		Items:      item.NewStore(),
	}
}
