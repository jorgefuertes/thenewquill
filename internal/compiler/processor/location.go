package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func readLocation(l line.Line, st *status.Status, a *adventure.Adventure) error {
	loc := location.New("", "")

	// continue reading location definition
	if st.HasCurrentLabel() {
		if !st.GetCurrentStoreable(&loc) {
			return cerr.ErrNoCurrentEntity.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename())
		}

		desc, ok := l.AsLocationDescription()
		if ok {
			loc.Description = desc
			st.CurrentStoreable = loc

			return nil
		}

		title, ok := l.AsLocationTitle()
		if ok {
			loc.Title = title
			st.CurrentStoreable = loc

			return nil
		}

		exitMap, ok := l.AsLocationConns()
		if ok {
			for action, dest := range exitMap {
				actionLabel, err := a.DB.AddLabel(action, false)
				if err != nil {
					return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
						WithFilename(st.CurrentFilename()).AddErr(err)
				}

				destLabel, err := a.DB.AddLabel(dest, false)
				if err != nil {
					return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
						WithFilename(st.CurrentFilename()).AddErr(err)
				}

				loc.SetConn(actionLabel.ID, destLabel.ID)
			}

			st.CurrentStoreable = loc

			return nil
		}

		// var
		if err := tryReadEntityVar(l, st, a); err != nil {
			return err
		}
	}

	// new location
	labelName, ok := l.AsLocationLabel()
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
		st.CurrentStoreable = loc

		return nil
	}

	return cerr.ErrWrongExitsDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
		WithFilename(st.CurrentFilename())
}
