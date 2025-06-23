package processor

import (
	"strconv"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/item"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/rg"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func readItem(l line.Line, st *status.Status, a *adventure.Adventure) error {
	i := item.New(db.UndefinedLabel.ID, db.UndefinedLabel.ID)

	if st.HasCurrentLabel() {

		if !st.GetCurrentStoreable(&i) {
			return cerr.ErrNoCurrentEntity.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename())
		}

		desc, ok := l.GetTextForFirstFoundLabelName("description", "desc")
		if ok {
			i.Description = desc
			st.CurrentStoreable = i

			return nil
		}

		o := l.OptimizedText()

		if o == "is wearable" {
			i.Wearable = true
			st.CurrentStoreable = i

			return nil
		}

		if o == "is created" {
			i.Created = true
			st.CurrentStoreable = i

			return nil
		}

		if o == "is container" {
			i.Container = true
			st.CurrentStoreable = i

			return nil
		}

		// is at, is in, is worn by
		if rg.ItemAt.MatchString(o) {
			parts := rg.ItemAt.FindStringSubmatch(o)

			atLabel, err := a.DB.AddLabel(parts[2], false)
			if err != nil {
				return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename()).AddErr(err)
			}

			i.At = atLabel.ID
			i.Worn = strings.Contains(parts[1], "worn")

			st.CurrentStoreable = i

			return nil
		}

		// has weight, weights
		if rg.ItemWeight.MatchString(o) {
			parts := rg.ItemWeight.FindStringSubmatch(o)

			var err error
			i.Weight, err = strconv.Atoi(parts[1])
			if err != nil {
				return cerr.ErrInvalidNumberDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename()).AddErr(err)
			}

			st.CurrentStoreable = i

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

			st.CurrentStoreable = i

			return nil
		}

		// var
		if err := tryReadEntityVar(l, st, a); err != nil {
			return err
		}
	}

	// new item
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

		nounLabel, err := a.DB.AddLabel(nounName, false)
		if err != nil {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		i.NounID = nounLabel.ID

		adjLabel, err := a.DB.AddLabel(adjName, false)
		if err != nil {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		i.AdjectiveID = adjLabel.ID

		st.CurrentStoreable = i

		return nil
	}

	return cerr.ErrWrongItemDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
		WithFilename(st.CurrentFilename())
}
