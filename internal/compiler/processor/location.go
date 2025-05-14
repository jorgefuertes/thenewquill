package processor

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/words"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/rg"
	"github.com/jorgefuertes/thenewquill/internal/compiler/section"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

func readLocation(l line.Line, st *status.Status, a *adventure.Adventure) error {
	label, ok := l.AsLocationLabel()
	if ok {
		st.CurrentLabel = label

		_, err := a.Locations.New(label)
		if err != nil {
			return cerr.ErrCannotCreateLocation.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename()).AddErr(err)
		}

		st.SetDef(label, section.Locs)

		return nil
	}

	if !st.HasCurrentLabel() {
		return cerr.ErrWrongLocationLabelDeclaration.WithSection(st.Section).WithStack(st.Stack).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	currentLocation := a.Locations.Get(st.CurrentLabel)
	if currentLocation == nil {
		return cerr.ErrWrongLocationLabelDeclaration.WithSection(st.Section).WithStack(st.Stack).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	o := l.OptimizedText()

	// vars
	if rg.Var.MatchString(o) {
		m := rg.Var.FindStringSubmatch(o)

		if !rg.IsValidLabel(m[1]) {
			return cerr.ErrInvalidLabel.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename())
		}

		currentLocation.Vars.SetFromString(m[1], m[2])

		return nil
	}

	desc, ok := l.AsLocationDescription()
	if ok {
		currentLocation.Description = desc

		return nil
	}

	title, ok := l.AsLocationTitle()
	if ok {
		currentLocation.Title = title

		return nil
	}

	exitMap, ok := l.AsLocationConns()
	if ok {
		for wordLabel, destLabel := range exitMap {
			word := a.Words.FirstWithTypes(wordLabel, words.Verb, words.Noun)
			if word == nil {
				word = a.Words.Set(wordLabel, words.Unknown)
				st.SetUndef(wordLabel, section.Words, l)
			}

			dest := a.Locations.Get(destLabel)
			if dest == nil {
				var err error
				dest, err = a.Locations.New(destLabel)
				if err != nil {
					return cerr.ErrCannotCreateLocation.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
						WithFilename(st.CurrentFilename()).AddErr(err)
				}

				st.SetUndef(destLabel, section.Locs, l)
			}

			currentLocation.SetConn(word, dest)
		}

		return nil
	}

	return cerr.ErrWrongExitsDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
		WithFilename(st.CurrentFilename())
}
