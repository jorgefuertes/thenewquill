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
		currentLabel, err := a.DB.GetLabel(st.GetCurrentLabelID())
		if err != nil {
			return true, cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		composedLabel := currentLabel + "." + name
		if _, err := a.Variables.SetByLabel(composedLabel, value); err != nil {
			return true, cerr.ErrWrongVariableDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		return true, nil
	}

	return false, nil
}

func readVar(l line.Line, st *status.Status, a *adventure.Adventure) error {
	label, value, ok := l.AsVar()
	if !ok {
		return cerr.ErrWrongVariableDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	if _, err := a.Variables.SetByLabel(label, value); err != nil {
		return cerr.ErrWrongVariableDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename()).AddErr(err)
	}

	return nil
}
