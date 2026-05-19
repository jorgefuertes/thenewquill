package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/message"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func readMessage(l line.Line, st *status.Status, a *adventure.Adventure) error {
	labelName, text, plural, ok := l.AsMsg()
	if !ok {
		return cerr.ErrWrongMessageDeclaration.WithSection(st.Section).WithStack(st.Stack).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	labelID, err := a.DB.CreateLabel(labelName)
	if err != nil {
		return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
			WithFilename(st.CurrentFilename()).AddErr(err)
	}

	// When the current entity is another plural variant of the same message,
	// fold this line into it instead of starting a new Message.
	if st.HasCurrent() && st.CurrentKind() == kind.Message && st.GetCurrentLabelID() == labelID {
		var current *message.Message
		if st.GetCurrentStoreable(&current) {
			current.SetPlural(plural, text)
			st.SetCurrentStoreable(current)

			return nil
		}
	}

	// New message: persist any pending entity from a previous section/label
	// and start a fresh one.
	if err := st.SaveCurrentStoreable(); !err.IsOK() {
		return err
	}

	m := message.New()
	m.LabelID = labelID
	m.SetPlural(plural, text)
	st.SetCurrentStoreable(m)

	return nil
}
