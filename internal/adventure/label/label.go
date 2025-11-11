package label

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
)

type Label struct {
	ID   id.ID
	Name string
}

var _ adapter.Storeable = &Label{}

var (
	Undefined  = Label{ID: id.Undefined, Name: "undefined"}
	Config     = Label{ID: 1, Name: "adv-config"}
	Underscore = Label{ID: 2, Name: "_"}
	Wildcard   = Label{ID: 3, Name: "*"}
)

func (l Label) GetID() id.ID {
	return l.ID
}

func (l Label) SetID(i id.ID) adapter.Storeable {
	l.ID = i

	return l
}
