package processor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/item"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/rg"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func readItem(l line.Line, st *status.Status, a *adventure.Adventure) error {
	if st.HasCurrent() {
		i := item.New()

		if !st.GetCurrentStoreable(&i) {
			return cerr.ErrNoCurrentEntity.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename())
		}

		desc, ok := l.GetTextForFirstFoundLabel("description", "desc")
		if ok {
			i.Description = desc
			st.SetCurrentStoreable(i)

			return nil
		}

		o := l.OptimizedText()

		if o == "is wearable" {
			i.Wearable = true
			st.SetCurrentStoreable(i)

			return nil
		}

		if o == "is created" {
			i.Created = true
			st.SetCurrentStoreable(i)

			return nil
		}

		if o == "is container" {
			i.Container = true
			st.SetCurrentStoreable(i)

			return nil
		}

		// is at, is in, is worn by
		if rg.ItemAt.MatchString(o) {
			parts := rg.ItemAt.FindStringSubmatch(o)

			atLabelID, err := a.DB.CreateLabel(parts[2])
			if err != nil {
				return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename()).AddErr(err)
			}

			i.At = atLabelID
			i.Worn = strings.Contains(parts[1], "worn")
			st.SetCurrentStoreable(i)

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
			st.SetCurrentStoreable(i)

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
			st.SetCurrentStoreable(i)

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
	label, nounLabel, adjLabel, ok := l.AsLabelNounAdjDeclaration()
	if ok {
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
				WithFilename(st.CurrentFilename()).AddErr(fmt.Errorf("%w while looking for noun %q", err, nounLabel))
		}

		adj, err := a.Words.GetAnyWith(adjLabel, word.Adjective)
		if err != nil {
			return cerr.ErrAdjectiveNotFound.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).
				AddErr(fmt.Errorf("%w while looking for adjective %q", err, adjLabel))
		}

		i := item.New()
		i.LabelID = labelID
		i.NounID = noun.ID
		i.AdjectiveID = adj.ID
		st.SetCurrentStoreable(i)

		return nil
	}

	return cerr.ErrWrongItemDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
		WithFilename(st.CurrentFilename())
}
