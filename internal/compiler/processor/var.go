package processor

import (
	"thenewquill/internal/adventure"
	cerr "thenewquill/internal/compiler/compiler_error"
	"thenewquill/internal/compiler/line"
	"thenewquill/internal/compiler/section"
	"thenewquill/internal/compiler/status"
)

func readVar(l line.Line, st *status.Status, a *adventure.Adventure) error {
	name, value, ok := l.AsVar()
	if !ok {
		return cerr.ErrWrongVariableDeclaration.WithStack(st.Stack).WithLine(l).WithFilename(st.CurrentFilename())
	}

	a.Vars.Set(name, value)
	st.SetDef(name, section.Vars)

	return nil
}
