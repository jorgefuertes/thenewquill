package processor

import (
	"strings"

	"thenewquill/internal/adventure"
	"thenewquill/internal/adventure/character"
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/words"
	cerr "thenewquill/internal/compiler/compiler_error"
	"thenewquill/internal/compiler/line"
	"thenewquill/internal/compiler/rg"
	"thenewquill/internal/compiler/section"
	"thenewquill/internal/compiler/status"
)

func readCharacter(l line.Line, st *status.Status, a *adventure.Adventure) error {
	if st.HasCurrentLabel() {
		c := a.Chars.Get(st.CurrentLabel)
		if c == nil {
			return cerr.ErrWrongCharDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename())
		}

		desc, ok := l.GetTextForFirstFoundLabel("description", "desc")
		if ok {
			c.Description = desc

			return nil
		}

		o := l.OptimizedText()

		if o == "is created" {
			c.Created = true

			return nil
		}

		if o == "is human" {
			h := a.Chars.GetHuman()
			if h != nil && h.Label != c.Label {
				return cerr.ErrOnlyOneHuman.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename())
			}

			c.Human = true

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
				inLoc = a.Locations.Set(locLabel, loc.Undefined, loc.Undefined)
				st.SetUndef(locLabel, section.Locs, l)
			}

			c.Location = inLoc

			return nil
		}

		if rg.Var.MatchString(o) {
			m := rg.Var.FindStringSubmatch(o)
			c.Vars.SetFromString(m[1], m[2])

			return nil
		}
	}

	label, noun, adj, ok := l.AsLabelNounAdjDeclaration()
	if ok {
		st.CurrentLabel = label
		st.SetDef(label, section.Chars)

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

		if err := a.Chars.Set(character.New(label, nounWord, adjWord)); err != nil {
			return cerr.ErrWrongCharDeclaration.WithStack(st.Stack).WithSection(st.Section).AddErr(err).WithLine(l).
				WithFilename(st.CurrentFilename())
		}

		st.SetDef(label, section.Chars)

		return nil
	}

	return cerr.ErrWrongCharDeclaration.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
		WithFilename(st.CurrentFilename())
}
