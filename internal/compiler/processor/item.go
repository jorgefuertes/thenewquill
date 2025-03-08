package processor

import (
	"strconv"

	"thenewquill/internal/adventure"
	"thenewquill/internal/adventure/item"
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/words"
	cerr "thenewquill/internal/compiler/compiler_error"
	"thenewquill/internal/compiler/line"
	"thenewquill/internal/compiler/rg"
	"thenewquill/internal/compiler/section"
	"thenewquill/internal/compiler/status"
)

func readItem(l line.Line, st *status.Status, a *adventure.Adventure) error {
	if st.HasCurrentLabel() {
		i := a.Items.Get(st.CurrentLabel)
		if i == nil {
			return cerr.ErrWrongItemDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename())
		}

		desc, ok := l.GetTextForFirstFoundLabel("description", "desc")
		if ok {
			i.Description = desc

			return nil
		}

		o := l.OptimizedText()

		switch o {
		case "is wearable":
			i.IsWearable = true

			return nil
		case "is worn":
			i.IsWearable = true
			i.Wear()

			return nil
		case "is created":
			i.Create()

			return nil
		case "is container":
			i.IsContainer = true

			return nil
		case "is held":
			i.Hold()

			return nil
		case "is destroyed", "is not created", "is not held", "is not worn", "is not wearable":
			return nil
		}

		if rg.IsAtLocation.MatchString(o) {
			m := rg.IsAtLocation.FindStringSubmatch(o)

			inLoc := a.Locations.Get(m[1])
			if inLoc == nil {
				inLoc = a.Locations.Set(m[1], loc.Undefined, loc.Undefined)
				st.SetUndef(m[1], section.Locs, l)
			}

			i.Location = inLoc

			return nil
		}

		if rg.ItemWeight.MatchString(o) {
			m := rg.ItemWeight.FindStringSubmatch(o)
			w, err := strconv.Atoi(m[1])
			if err != nil {
				return cerr.ErrWrongItemWeight.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename())
			}

			i.Weight = w

			return nil
		}

		if rg.ItemMaxWeight.MatchString(o) {
			m := rg.ItemMaxWeight.FindStringSubmatch(o)
			w, err := strconv.Atoi(m[1])
			if err != nil {
				return cerr.ErrWrongItemWeight.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename())
			}

			i.MaxWeight = w

			return nil
		}

		if rg.Var.MatchString(o) {
			m := rg.Var.FindStringSubmatch(o)
			i.Vars.Set(m[1], m[2])

			return nil
		}
	}

	label, noun, adj, ok := l.AsLabelNounAdjDeclaration()
	if ok {
		if a.Items.Exists(label) {
			return cerr.ErrDuplicatedItemLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename())
		}

		st.CurrentLabel = label
		st.SetDef(label, section.Items)

		nounWord := a.Words.Get(words.Noun, noun)
		if nounWord == nil {
			nounWord = a.Words.Set(noun, words.Noun)
			st.SetUndef(noun, section.Words, l)
		}

		adjWord := a.Words.Get(words.Adjective, adj)
		if adjWord == nil {
			adjWord = a.Words.Set(adj, words.Adjective)
			st.SetUndef(adj, section.Words, l)
		}

		if err := a.Items.Set(item.New(label, nounWord, adjWord)); err != nil {
			return cerr.ErrWrongItemDeclaration.WithStack(st.Stack).WithSection(st.Section).AddErr(err).WithLine(l).
				WithFilename(st.CurrentFilename())
		}

		st.SetDef(label, section.Items)

		return nil
	}

	return cerr.ErrWrongItemDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
		WithFilename(st.CurrentFilename())
}
