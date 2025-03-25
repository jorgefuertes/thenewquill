package processor

import (
	"strconv"
	"strings"

	"thenewquill/internal/adventure"
	"thenewquill/internal/adventure/item"
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

		if o == "is wearable" {
			i.IsWearable = true

			return nil
		}

		if o == "is worn" {
			i.IsWearable = true
			i.Wear()

			return nil
		}

		if o == "is created" {
			i.Create()

			return nil
		}

		if o == "is container" {
			i.IsContainer = true

			return nil
		}

		if strings.HasPrefix(o, "is at ") {
			locLabel := strings.TrimPrefix(o, "is at ")
			if !rg.IsValidLabel(locLabel) {
				return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename())
			}

			inLoc := a.Locations.Get(locLabel)
			if inLoc == nil {
				var err error
				inLoc, err = a.Locations.New(locLabel)
				if err != nil {
					return cerr.ErrCannotCreateLocation.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
						WithFilename(st.CurrentFilename()).AddErr(err)
				}

				st.SetUndef(locLabel, section.Locs, l)
			}

			i.Location = inLoc

			return nil
		}

		if strings.HasPrefix(o, "is in ") {
			containerLabel := strings.TrimPrefix(o, "is in ")
			if !rg.IsValidLabel(containerLabel) {
				return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename())
			}

			container := a.Items.Get(containerLabel)
			if container == nil {
				container := item.New(containerLabel, nil, nil)

				if err := a.Items.Set(container); err != nil {
					return cerr.ErrWrongItemDeclaration.WithStack(st.Stack).
						WithSection(st.Section).
						AddErr(err).
						WithLine(l).
						WithFilename(st.CurrentFilename())
				}

				st.SetUndef(containerLabel, section.Items, l)
			}

			i.Inside = container

			return nil
		}

		if strings.HasPrefix(o, "has weight ") {
			w, err := strconv.Atoi(strings.TrimPrefix(o, "has weight "))
			if err != nil {
				return cerr.ErrWrongItemDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename())
			}

			i.Weight = w

			return nil
		}

		if strings.HasPrefix(o, "has max weight ") {
			w, err := strconv.Atoi(strings.TrimPrefix(o, "has max weight "))
			if err != nil {
				return cerr.ErrWrongItemDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename())
			}

			i.MaxWeight = w

			return nil
		}

		// vars
		if rg.Var.MatchString(o) {
			m := rg.Var.FindStringSubmatch(o)

			if !rg.IsValidLabel(m[1]) {
				return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename())
			}

			i.Vars.SetFromString(m[1], m[2])

			return nil
		}
	}

	label, noun, adj, ok := l.AsLabelNounAdjDeclaration()
	if ok {
		if !rg.IsValidLabel(label) {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
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
