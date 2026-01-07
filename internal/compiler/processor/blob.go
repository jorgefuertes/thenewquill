package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/blob"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

// readBlob processes a blob declaration line and adds it to the adventure.
func readBlob(l line.Line, st *status.Status, a *adventure.Adventure) error {
	label, filename, ok := l.AsBlob()
	if !ok {
		return cerr.ErrWrongBlobDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	labelID, err := a.DB.CreateLabel(label)
	if err != nil {
		return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename()).AddErr(err)
	}

	b := blob.New()
	b.SetLabelID(labelID)
	if err := b.Load(st.CurrentPath(filename)); err != nil {
		return cerr.ErrLoadingBlob.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename()).AddErr(err)
	}

	if _, err := a.Blobs.Create(b); err != nil {
		return cerr.ErrCannotCreateBlob.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename()).AddErr(err)
	}

	return nil
}
