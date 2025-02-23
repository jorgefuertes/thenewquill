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
		return cerr.ErrWrongWordDeclaration.WithStack(st.Stack).WithLine(l).WithFilename(st.CurrentFilename())
	}

	_ = a.Vocabulary.Set(w.Label, w.Type, w.Synonyms...)
	st.SetDef(w.Label, section.Words)

	return nil
}
