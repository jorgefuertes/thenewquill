package processor

import (
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func readCharacter(l line.Line, st *status.Status, a *adventure.Adventure) error {
	c := character.New(db.UndefinedLabel.ID, db.UndefinedLabel.ID)

	if st.HasCurrentLabel() {
		if !st.GetCurrentStoreable(&c) {
			return cerr.ErrNoCurrentEntity.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename())
		}

		desc, ok := l.GetTextForFirstFoundLabelName("description", "desc")
		if ok {
			c.Description = desc
			st.CurrentStoreable = c

			return nil
		}

		o := l.OptimizedText()

		if o == "is created" {
			c.Created = true
			st.CurrentStoreable = c

			return nil
		}

		if o == "is human" {
			if a.Characters.HasHuman() {
				return cerr.ErrOnlyOneHuman.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename())
			}

			c.Human = true
			st.CurrentStoreable = c

			return nil
		}

		if strings.HasPrefix(o, "is at ") {
			locName := strings.TrimPrefix(o, "is at ")

			locLabel, err := a.DB.AddLabel(locName, false)
			if err != nil {
				return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename()).AddErr(err)
			}

			c.LocationID = locLabel.ID
			st.CurrentStoreable = c

			return nil
		}

		// var
		if err := tryReadEntityVar(l, st, a); err != nil {
			return err
		}
	}

	labelName, nounName, adjName, ok := l.AsLabelNounAdjDeclaration()
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

		nounLabel, err := a.DB.GetLabelByName(nounName)
		if err != nil {
			nounLabel, err = a.DB.AddLabel(nounName, false)
			if err != nil {
				return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename()).AddErr(err)
			}
		}

		adjLabel, err := a.DB.GetLabelByName(adjName)
		if err != nil {
			adjLabel, err = a.DB.AddLabel(adjName, false)
			if err != nil {
				return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename()).AddErr(err)
			}
		}

		st.CurrentStoreable = character.New(nounLabel.ID, adjLabel.ID)

		return nil
	}

	return cerr.ErrWrongCharDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
		WithFilename(st.CurrentFilename())
}
