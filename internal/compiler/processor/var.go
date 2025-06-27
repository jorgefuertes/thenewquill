package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func tryReadEntityVar(l line.Line, st *status.Status, a *adventure.Adventure) (bool, error) {
	name, value, ok := l.AsVar()
	if ok {
		vLabel, err := a.DB.AddLabel(st.GetCurrentLabel().Name+"."+name, true)
		if err != nil {
			return true, cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		if err := a.Variables.Set(vLabel.ID, value); err != nil {
			return true, cerr.ErrWrongVariableDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		return true, nil
	}

	return false, nil
}

func readVar(l line.Line, st *status.Status, a *adventure.Adventure) error {
	name, value, ok := l.AsVar()
	if !ok {
		return cerr.ErrWrongVariableDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	vLabel, err := a.DB.AddLabel(name, false)
	if err != nil {
		return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename()).AddErr(err)
	}

	if err := a.Variables.Set(vLabel.ID, value); err != nil {
		return cerr.ErrWrongVariableDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename()).AddErr(err)
	}

	return nil
}
