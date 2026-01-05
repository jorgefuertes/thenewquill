package processor

import (
	"errors"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
	"github.com/jorgefuertes/thenewquill/internal/database"
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
			st.SetCurrentStoreable(c)

			return nil
		}

		o := l.OptimizedText()

		if o == "is created" {
			c.Created = true
			st.SetCurrentStoreable(c)

			return nil
		}

		if o == "is human" {
			if a.DB.Query(database.FilterByKind(kind.Character), database.NewFilter("Human", database.Equal, true)).
				Exists() {
				return cerr.ErrOnlyOneHuman.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename())
			}

			c.Human = true
			st.SetCurrentStoreable(c)

			return nil
		}

		if strings.HasPrefix(o, "is at ") {
			locLabel := strings.TrimPrefix(o, "is at ")
			loc, err := a.Locations.Get().WithLabel(locLabel).First()
			if err != nil {
				return cerr.ErrLocationNotFound.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename()).AddErr(err)
			}

			c.LocationID = loc.ID
			st.SetCurrentStoreable(c)

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

	label, nounLabel, adjLabel, ok := l.AsLabelNounAdjDeclaration()
	if ok {
		// close and save current character if any
		if err := st.SaveCurrentStoreable(); !err.IsOK() {
			return err
		}

		labelID, err := a.DB.CreateLabel(label)
		if err != nil {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		noun, err := a.Words.GetAnyWith(nounLabel, word.Noun)
		if err != nil {
			return cerr.ErrWordNotFound.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		adj, err := a.Words.GetAnyWith(adjLabel, word.Adjective)
		if err != nil {
			return cerr.ErrWordNotFound.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		st.SetCurrentStoreable(&character.Character{LabelID: labelID, NounID: noun.ID, AdjectiveID: adj.ID})

		return nil
	}

	return cerr.ErrWrongCharDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
		WithFilename(st.CurrentFilename())
}
