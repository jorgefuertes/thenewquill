package processor

import (
	"errors"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func readLocation(l line.Line, st *status.Status, a *adventure.Adventure) error {
	loc := location.New("", "")

	// continue reading location definition
	if st.HasCurrent() {
		if !st.GetCurrentStoreable(&loc) {
			return errors.New("unexpected: cannot get current location")
		}

		desc, ok := l.AsLocationDescription()
		if ok {
			loc.Description = desc

			if err := st.SetCurrentStoreable(loc); err != nil {
				return err
			}

			return nil
		}

		title, ok := l.AsLocationTitle()
		if ok {
			loc.Title = title

			if err := st.SetCurrentStoreable(loc); err != nil {
				return err
			}

			return nil
		}

		exitMap, ok := l.AsLocationConns()
		if ok {
			for action, dest := range exitMap {
				actionLabel, err := a.DB.AddLabel(action)
				if err != nil {
					return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
						WithFilename(st.CurrentFilename()).AddErr(err)
				}

				destLabel, err := a.DB.AddLabel(dest)
				if err != nil {
					return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
						WithFilename(st.CurrentFilename()).AddErr(err)
				}

				loc.SetConn(actionLabel.ID, destLabel.ID)
			}

			if err := st.SetCurrentStoreable(loc); err != nil {
				return err
			}

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
		if err := st.SaveCurrentStoreable(); !err.IsOK() {
			return err
		}

		labelID, _, err := a.DB.CreateLabelIfNotExists(labelName, false)
		if err != nil {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		if err := st.SetCurrentLabelID(labelID); err != nil {
			return err
		}

		if err := st.SetCurrentStoreable(loc); err != nil {
			return err
		}

		return nil
	}

	return cerr.ErrWrongLocationDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
		WithFilename(st.CurrentFilename())
}
