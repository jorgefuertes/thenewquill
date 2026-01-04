package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
	"github.com/jorgefuertes/thenewquill/pkg/log"
)

func ProcessLine(l line.Line, st *status.Status, a *adventure.Adventure) error {
	log.Debug("⚙︎ Processing line %q for section %q", l.Text, st.Section.String())

	switch st.Section {
	case kind.Param:
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
