package adventure

import (
	"fmt"

	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/msg"
	"thenewquill/internal/adventure/obj"
	"thenewquill/internal/adventure/store"
	"thenewquill/internal/adventure/voc"
)

type Adventure struct {
	Vars       store.Store
	Vocabulary voc.Vocabulary
	Messages   msg.Messages
	Locations  loc.Locations
	Objects    obj.Items
}

func New() *Adventure {
	return &Adventure{
		Vars:       store.New(),
		Vocabulary: voc.New(),
		Messages:   msg.New(),
		Locations:  loc.New(),
		Objects:    obj.NewItems(),
	}
}

func (a *Adventure) Dump() {
	for _, m := range a.Messages {
		fmt.Printf("[%s] %s: %s\n", m.Type, m.Label, m.Text)
	}
}
