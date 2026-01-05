package processor

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func readLocation(l line.Line, st *status.Status, a *adventure.Adventure) error {
	// continue reading location definition if one present
	if st.HasCurrent() {
		loc := location.New()

		if !st.GetCurrentStoreable(&loc) {
			panic("unexpected: cannot get current location")
		}

		desc, ok := l.AsLocationDescription()
		if ok {
			loc.Description = desc
			st.SetCurrentStoreable(loc)

			return nil
		}

		title, ok := l.AsLocationTitle()
		if ok {
			loc.Title = title
			st.SetCurrentStoreable(loc)

			return nil
		}

		exitMap, ok := l.AsLocationConns()
		if ok {
			for actionLabel, destLabel := range exitMap {
				actionWord, err := a.Words.GetAnyWith(actionLabel, word.Verb, word.Noun)
				if err != nil {
					return cerr.ErrWordNotFound.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
						WithFilename(st.CurrentFilename()).
						AddErr(err).AddErr(fmt.Errorf("missing actionWord with label %q", actionLabel))
				}

				// assign a label temporarily to the destination
				// real ID will be assigned later
				destLabelID, err := a.DB.CreateLabel(destLabel)
				if err != nil {
					return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
						WithFilename(st.CurrentFilename()).AddErr(err)
				}

				loc.SetConn(actionWord.ID, destLabelID)
			}

			st.SetCurrentStoreable(loc)

			return nil
		}

		// var
		isVar, err := tryReadEntityVar(l, st, a)
		if err != nil {
			return err
		}

		if isVar {
			return nil
		}
	}

	// new location
	labelName, ok := l.AsLocationLabel()
	if ok {
		if st.HasCurrent() {
			if err := st.SaveCurrentStoreable(); !err.IsOK() {
				return err
			}
		}

		labelID, err := a.DB.CreateLabel(labelName)
		if err != nil {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		st.SetCurrentStoreable(&location.Location{LabelID: labelID})

		return nil
	}

	return cerr.ErrWrongLocationDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
		WithFilename(st.CurrentFilename())
}
