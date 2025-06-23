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
		// save current storeable if any
		if err := st.Save(a.DB); err != nil {
			return cerr.ErrDBCreate.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		label, err := a.DB.AddLabel(labelName, false)
		if err != nil {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		st.CurrentLabel = label
		m.SetPlural(plural, text)
		st.CurrentStoreable = m

		return nil
	}

	return cerr.ErrWrongMessageDeclaration.WithSection(st.Section).WithStack(st.Stack).WithLine(l).
		WithFilename(st.CurrentFilename())
}
