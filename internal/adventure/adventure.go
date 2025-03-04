package adventure

import (
	"thenewquill/internal/adventure/character"
	"thenewquill/internal/adventure/config"
	"thenewquill/internal/adventure/item"
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/msg"
	"thenewquill/internal/adventure/vars"
	"thenewquill/internal/adventure/words"
)

type Adventure struct {
	Config    config.Config
	Vars      vars.Store
	Words     words.Store
	Messages  msg.Store
	Locations loc.Store
	Items     item.Store
	Chars     character.Store
}

func New() *Adventure {
	return &Adventure{
		Config:    config.New(),
		Vars:      vars.NewStore(),
		Words:     words.NewStore(),
		Messages:  msg.NewStore(),
		Locations: loc.NewStore(),
		Items:     item.NewStore(),
		Chars:     character.NewStore(),
	}
}
