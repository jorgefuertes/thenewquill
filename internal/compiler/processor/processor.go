package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func ProcessLine(l line.Line, st *status.Status, a *adventure.Adventure) error {
	switch st.Section {
	case kind.Config:
		return readConfig(l, st, a)
	case kind.Variable:
		return readVar(l, st, a)
	case kind.Word:
		return readWord(l, st, a)
	case kind.Message:
		return readMessage(l, st, a)
	case kind.Item:
		return readItem(l, st, a)
	case kind.Location:
		return readLocation(l, st, a)
	case kind.Character:
		return readCharacter(l, st, a)
	default:
		return nil
	}
}
