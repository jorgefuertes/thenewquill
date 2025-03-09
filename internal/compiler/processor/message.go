package processor

import (
	"thenewquill/internal/adventure"
	"thenewquill/internal/adventure/msg"
	cerr "thenewquill/internal/compiler/compiler_error"
	"thenewquill/internal/compiler/line"
	"thenewquill/internal/compiler/status"
)

func readMessage(l line.Line, st *status.Status, a *adventure.Adventure) error {
	m, ok := l.AsMsg()
	if !ok {
		return cerr.ErrWrongMessageDeclaration.WithSection(st.Section).WithStack(st.Stack).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	if err := a.Messages.Set(m); err != nil {
		return cerr.ErrWrongMessageDeclaration.WithStack(st.Stack).WithSection(st.Section).AddErr(err).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	st.SetDef(m.Label, st.Section)

	if m.IsPluralized() {
		// recover it from the store
		m = a.Messages.Get(m.Label)
		if m == nil {
			return cerr.ErrCannotRetrieveMessage.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename())
		}

		// define and undefine reminders for the plurals
		for i, text := range m.Plurals {
			pLabel := m.Label + "." + msg.PluralNames[i]
			if text == "" && !st.IsUndef(pLabel, st.Section) {
				st.SetUndef(pLabel, st.Section, l)

				continue
			}
			st.SetDef(pLabel, st.Section)
		}
	}

	return nil
}
