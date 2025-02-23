package processor

import (
	"thenewquill/internal/adventure"
	"thenewquill/internal/adventure/msg"
	cerr "thenewquill/internal/compiler/compiler_error"
	"thenewquill/internal/compiler/line"
	"thenewquill/internal/compiler/section"
	"thenewquill/internal/compiler/status"
)

func readMessage(l line.Line, st *status.Status, a *adventure.Adventure) error {
	t := msg.SystemMsg
	if st.Section == section.UserMsgs {
		t = msg.UserMsg
	}

	m, ok := l.AsMsg(t)
	if !ok {
		return cerr.ErrWrongMessageDeclaration.WithStack(st.Stack).WithLine(l).WithFilename(st.CurrentFilename())
	}

	if err := a.Messages.Set(m); err != nil {
		return cerr.ErrWrongMessageDeclaration.WithStack(st.Stack).
			AddErr(err).
			WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	st.SetDef(m.Label, section.SysMsg)

	return nil
}
