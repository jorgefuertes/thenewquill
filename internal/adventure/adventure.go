package adventure

import (
	"fmt"

	"thenewquill/internal/adventure/character"
	"thenewquill/internal/adventure/config"
	"thenewquill/internal/adventure/item"
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/msg"
	"thenewquill/internal/adventure/vars"
	"thenewquill/internal/adventure/words"
	"thenewquill/internal/compiler/section"
)

const VERSION = 1

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

func (a *Adventure) ExportHeaders() []string {
	return []string{
		fmt.Sprintf("Database version %03d with %03d sections", VERSION, len(a.Export())),
		fmt.Sprintf("%s v%s", a.Config.Title, a.Config.Version),
		fmt.Sprintf("By %s", a.Config.Author),
	}
}

func (a *Adventure) Export() map[section.Section][][]string {
	exportFuncs := []func() (section.Section, [][]string){
		a.Config.Export,
		a.Vars.Export,
		a.Words.Export,
		a.Messages.Export,
		a.Locations.Export,
		a.Items.Export,
		a.Chars.Export,
	}

	data := make(map[section.Section][][]string, 0)
	for _, f := range exportFuncs {
		sec, rows := f()
		data[sec] = rows
	}

	return data
}
