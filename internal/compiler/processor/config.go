package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
	"github.com/jorgefuertes/thenewquill/pkg/log"
)

func readConfig(l line.Line, st *status.Status, a *adventure.Adventure) error {
	field, value, ok := l.AsConfig()
	if !ok {
		return cerr.ErrWrongConfigDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	log.Debug("💾 [readConfig] SetConfig: %s = %s", field, value)

	id, err := a.Config.Set(field, value)
	if err == nil {
		log.Debug("🗄️ [DB] Created param %q with ID %d", field, id)
	}

	return err
}
