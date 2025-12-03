package processor

import (
	"strconv"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/item"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/rg"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func readItem(l line.Line, st *status.Status, a *adventure.Adventure) error {
	i := item.New(primitive.UndefinedID, primitive.UndefinedID, primitive.UndefinedID)

	if st.HasCurrent() {
		if !st.GetCurrentStoreable(&i) {
			return cerr.ErrNoCurrentEntity.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename())
		}

		desc, ok := l.GetTextForFirstFoundLabel("description", "desc")
		if ok {
			i.Description = desc

			if err := st.SetCurrentStoreable(i); err != nil {
				return err
			}

			return nil
		}

		o := l.OptimizedText()

		if o == "is wearable" {
			i.Wearable = true

			if err := st.SetCurrentStoreable(i); err != nil {
				return err
			}

			return nil
		}

		if o == "is created" {
			i.Created = true

			if err := st.SetCurrentStoreable(i); err != nil {
				return err
			}

			return nil
		}

		if o == "is container" {
			i.Container = true

			if err := st.SetCurrentStoreable(i); err != nil {
				return err
			}

			return nil
		}

		// is at, is in, is worn by
		if rg.ItemAt.MatchString(o) {
			parts := rg.ItemAt.FindStringSubmatch(o)

			atLabelID, _, err := a.DB.CreateLabelFromString(parts[2], false)
			if err != nil {
				return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename()).AddErr(err)
			}

			i.At = atLabelID
			i.Worn = strings.Contains(parts[1], "worn")

			if err := st.SetCurrentStoreable(i); err != nil {
				return err
			}

			return nil
		}

		// has weight, weights
		if rg.ItemWeight.MatchString(o) {
			parts := rg.ItemWeight.FindStringSubmatch(o)

			var err error
			i.Weight, err = strconv.Atoi(parts[2])
			if err != nil {
				return cerr.ErrInvalidNumberDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename()).AddErr(err)
			}

			if err := st.SetCurrentStoreable(i); err != nil {
				return err
			}

			return nil
		}

		// has max weight, max weight
		if rg.ItemMaxWeight.MatchString(o) {
			parts := rg.ItemMaxWeight.FindStringSubmatch(o)

			var err error
			i.MaxWeight, err = strconv.Atoi(parts[1])
			if err != nil {
				return cerr.ErrInvalidNumberDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename()).AddErr(err)
			}

			if err := st.SetCurrentStoreable(i); err != nil {
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

	// new item
	labelName, nounName, adjName, ok := l.AsLabelNounAdjDeclaration()
	if ok {
		if err := st.SaveCurrentStoreable(); !err.IsOK() {
			return err
		}

		labelID, _, err := a.DB.CreateLabelFromString(labelName, false)
		if err != nil {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		if err := st.SetCurrentLabelID(labelID); err != nil {
			return err
		}

		nounLabelID, _, err := a.DB.CreateLabelFromString(nounName, false)
		if err != nil {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		i.NounID = nounLabelID

		adjLabelID, _, err := a.DB.CreateLabelFromString(adjName, false)
		if err != nil {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		i.AdjectiveID = adjLabelID

		if err := st.SetCurrentStoreable(i); err != nil {
			return err
		}

		return nil
	}

	return cerr.ErrWrongItemDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
		WithFilename(st.CurrentFilename())
}
