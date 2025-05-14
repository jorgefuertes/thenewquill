package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func readConfig(l line.Line, st *status.Status, a *adventure.Adventure) error {
	label, value, ok := l.AsConfig()
	if !ok || label == config.UnknownField {
		return cerr.ErrWrongConfigDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	return a.Config.Set(label, value)
}
