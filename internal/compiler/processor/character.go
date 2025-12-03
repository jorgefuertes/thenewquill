package processor

import (
	"errors"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
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
			labelID, _, err := a.DB.CreateLabelFromString(locName, false)
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

		labelID, _, err := a.DB.CreateLabelFromString(labelName, false)
		if err := st.SetCurrentLabelID(labelID); err != nil {
			return err
		}

		nounLabelID, _, err := a.DB.CreateLabelFromString(nounName, false)
		if err != nil {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		adjLabelID, _, err := a.DB.CreateLabelFromString(adjName, false)
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
