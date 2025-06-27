package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/message"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func readMessage(l line.Line, st *status.Status, a *adventure.Adventure) error {
	m := message.New("")

	labelName, text, plural, ok := l.AsMsg()
	if ok {
		if err := st.SaveCurrentStoreable(); !err.IsOK() {
			return err
		}

		label, err := a.DB.AddLabel(labelName, false)
		if err != nil {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		if err := st.SetCurrentLabel(label); err != nil {
			return err
		}

		m.SetPlural(plural, text)

		if err := st.SetCurrentStoreable(m); err != nil {
			return err
		}

		return nil
	}

	return cerr.ErrWrongMessageDeclaration.WithSection(st.Section).WithStack(st.Stack).WithLine(l).
		WithFilename(st.CurrentFilename())
}
