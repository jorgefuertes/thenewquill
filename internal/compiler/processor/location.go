package processor

import (
	"thenewquill/internal/adventure"
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/voc"
	cerr "thenewquill/internal/compiler/compiler_error"
	"thenewquill/internal/compiler/line"
	"thenewquill/internal/compiler/section"
	"thenewquill/internal/compiler/status"
)

func readLocation(l line.Line, st *status.Status, a *adventure.Adventure) error {
	label, ok := l.AsLocationLabel()
	if ok {
		st.CurrentLabel = label
		a.Locations.Set(label, loc.Undefined, loc.Undefined)
		st.SetDef(label, section.Locs)

		return nil
	}

	if !st.HasCurrentLabel() {
		return cerr.ErrWrongLocationLabelDeclaration.WithSection(st.Section).WithStack(st.Stack).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	currentLocation := a.Locations.Get(st.CurrentLabel)

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
			word := a.Vocabulary.FirstWithTypes(wordLabel, voc.Verb, voc.Noun)
			if word == nil {
				word = a.Vocabulary.Set(wordLabel, voc.Unknown)
				st.SetUndef(wordLabel, section.Words, l)
			}

			dest := a.Locations.Get(destLabel)
			if dest == nil {
				dest = a.Locations.Set(destLabel, loc.Undefined, loc.Undefined)
				st.SetUndef(destLabel, section.Locs, l)
			}

			currentLocation.SetConn(word, dest)
		}

		return nil
	}

	return cerr.ErrWrongExitsDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
		WithFilename(st.CurrentFilename())
}
