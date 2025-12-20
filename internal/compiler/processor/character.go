package processor

import (
	"errors"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/internal/database/primitive"
)

func readCharacter(l line.Line, st *status.Status, a *adventure.Adventure) error {
	c := &character.Character{}

	if st.HasCurrent() {
		if !st.GetCurrentStoreable(&c) {
			return errors.New("unexpected: cannot get current character")
		}

		desc, ok := l.GetTextForFirstFoundLabel("description", "desc")
		if ok {
			c.Description = desc

			if err := st.SetCurrentStoreable(c); err != nil {
				return err
			}

			return nil
		}

		o := l.OptimizedText()

		if o == "is created" {
			c.Created = true

			if err := st.SetCurrentStoreable(c); err != nil {
				return err
			}

			return nil
		}

		if o == "is human" {
			if a.Characters.HasHuman() {
				return cerr.ErrOnlyOneHuman.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename())
			}

			c.Human = true

			if err := st.SetCurrentStoreable(c); err != nil {
				return err
			}

			return nil
		}

		if strings.HasPrefix(o, "is at ") {
			locName := strings.TrimPrefix(o, "is at ")

			labelID, err := a.DB.CreateLabelIfNotExists(locName, database.DenyCompositeLabel)
			if err != nil {
				return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename()).AddErr(err)
			}

			c.LocationID = labelID

			if err := st.SetCurrentStoreable(c); err != nil {
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

	labelName, nounName, adjName, ok := l.AsLabelNounAdjDeclaration()
	if ok {
		if err := st.SaveCurrentStoreable(); !err.IsOK() {
			return err
		}

		labelID, err := a.DB.CreateLabelIfNotExists(labelName, database.DenyCompositeLabel)
		if err := st.SetCurrentLabelID(labelID); err != nil {
			return err
		}

		nounLabelID, err := a.DB.CreateLabelIfNotExists(nounName, database.DenyCompositeLabel)
		if err != nil {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		adjLabelID, err := a.DB.CreateLabelIfNotExists(adjName, database.DenyCompositeLabel)
		if err != nil {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		c := character.New(primitive.UndefinedID, labelID, nounLabelID, adjLabelID)
		if err := st.SetCurrentStoreable(c); err != nil {
			return err
		}

		return nil
	}

	return cerr.ErrWrongCharDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
		WithFilename(st.CurrentFilename())
}
