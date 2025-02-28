package processor

import (
	"thenewquill/internal/adventure"
	"thenewquill/internal/compiler/line"
	"thenewquill/internal/compiler/section"
	"thenewquill/internal/compiler/status"
)

func ProcessLine(l line.Line, st *status.Status, a *adventure.Adventure) error {
	switch st.Section {
	case section.Config:
		return readConfig(l, st, a)
	case section.Vars:
		return readVar(l, st, a)
	case section.Words:
		return readWord(l, st, a)
	case section.SysMsg, section.UserMsgs:
		return readMessage(l, st, a)
	case section.Items:
		return readItem(l, st, a)
	case section.Locs:
		return readLocation(l, st, a)
	default:
		return nil
	}
}
