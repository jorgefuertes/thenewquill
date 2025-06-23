package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func readWord(l line.Line, st *status.Status, a *adventure.Adventure) error {
	kind, syns, ok := l.AsWord()
	if !ok {
		return cerr.ErrWrongWordDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	wordType := word.WordTypeFromString(kind)
	if wordType == word.None {
		return cerr.ErrWrongWordDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	if len(syns) == 0 {
		return cerr.ErrWrongWordDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	label, err := a.DB.AddLabel(syns[0], false)
	if err != nil {
		return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename()).AddErr(err)
	}

	w := word.New(wordType, syns...)
	w.ID = label.ID

	if err := a.Words.Create(w); err != nil {
		return cerr.ErrDBCreate.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename()).AddErr(err)
	}

	return nil
}
