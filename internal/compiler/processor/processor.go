package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/section"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func ProcessLine(l line.Line, st *status.Status, a *adventure.Adventure) error {
	switch st.Section {
	case section.Config:
		return readConfig(l, st, a)
	case section.Vars:
		return readVar(l, st, a)
	case section.Words:
		return readWord(l, st, a)
	case section.Messages:
		return readMessage(l, st, a)
	case section.Items:
		return readItem(l, st, a)
	case section.Locs:
		return readLocation(l, st, a)
	case section.Chars:
		return readCharacter(l, st, a)
	default:
		return nil
	}
}
