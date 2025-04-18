package processor

import (
	"thenewquill/internal/adventure"
	cerr "thenewquill/internal/compiler/compiler_error"
	"thenewquill/internal/compiler/line"
	"thenewquill/internal/compiler/section"
	"thenewquill/internal/compiler/status"
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
