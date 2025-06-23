package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func ProcessLine(l line.Line, st *status.Status, a *adventure.Adventure) error {
	switch st.Section {
	case db.Config:
		return readConfig(l, st, a)
	case db.Variables:
		return readVar(l, st, a)
	case db.Words:
		return readWord(l, st, a)
	case db.Messages:
		return readMessage(l, st, a)
	case db.Items:
		return readItem(l, st, a)
	case db.Locations:
		return readLocation(l, st, a)
	case db.Characters:
		return readCharacter(l, st, a)
	default:
		return nil
	}
}
