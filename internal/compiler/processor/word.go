package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/section"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func readWord(l line.Line, st *status.Status, a *adventure.Adventure) error {
	w, ok := l.AsWord()
	if !ok {
		return cerr.ErrWrongWordDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	for _, syn := range w.Synonyms {
		if a.Words.Exists(w.Type, syn) {
			existent := a.Words.Get(w.Type, syn)
			return cerr.ErrDuplicatedSynonym.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).
				AddMsgf("synonym '%s' present in %s '%s'", syn, w.Type.String(), existent.Label)
		}
	}

	_ = a.Words.Set(w.Label, w.Type, w.Synonyms...)
	st.SetDef(w.Label, section.Words)

	return nil
}
