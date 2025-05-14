package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/section"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func readVar(l line.Line, st *status.Status, a *adventure.Adventure) error {
	name, value, ok := l.AsVar()
	if !ok {
		return cerr.ErrWrongVariableDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	a.Vars.Set(name, value)
	st.SetDef(name, section.Vars)

	return nil
}
