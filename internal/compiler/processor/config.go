package processor

import (
	"thenewquill/internal/adventure"
	"thenewquill/internal/adventure/config"
	cerr "thenewquill/internal/compiler/compiler_error"
	"thenewquill/internal/compiler/line"
	"thenewquill/internal/compiler/status"
)

func readConfig(l line.Line, st *status.Status, a *adventure.Adventure) error {
	label, value, ok := l.AsConfig()
	if !ok || label == config.UnknownLabel {
		return cerr.ErrWrongConfigDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	return a.Config.Set(label, value)
}
