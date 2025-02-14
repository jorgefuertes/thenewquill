package compiler

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"thenewquill/internal/adventure"
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/msg"
	"thenewquill/internal/adventure/obj"
	"thenewquill/internal/adventure/voc"
)

func Compile(filename string) (*adventure.Adventure, error) {
	a := adventure.New()
	st := newStatus()

	err := compileFile(st, filename, a)
	cErr, ok := err.(compilerError)
	if ok {
		fmt.Println(cErr.Dump())
		return a, cErr
	}

	// check for unresolved labels
	for _, udf := range st.undef {
		fmt.Println(
			ErrUnresolvedLabel.WithFilename(udf.file).
				WithLine(udf.line).
				AddMsgf("%s `%s` remains undefined", udf.section.singleString(), udf.label).
				Dump(),
		)
		err = ErrRemainingUnresolvedLabels
	}

	return a, err
}

func compileFile(st *status, filename string, a *adventure.Adventure) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	n := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		n++
		l := newLine(scanner.Text(), n)

		if l.isBlank() {
			continue
		}

		st.appendStack(l)

		if l.isOneLineComment() {
			continue
		}

		// comments are ignored
		if l.isCommentBegin() && !st.comment.isOn() {
			st.setComment(l)

			continue
		}

		if st.comment.isOn() && l.isCommentEnd() {
			st.unsetComment()

			continue
		}

		if st.comment.isOn() {
			continue
		}

		// follow includes
		f, ok := l.toInClude()
		if ok {
			err := compileFile(st, filepath.Dir(filename)+"/"+f, a)
			if err != nil {
				cErr, ok := err.(compilerError)
				if ok {
					return cErr
				}

				return ErrCannotOpenIncludedFile.WithStack(st.stack).WithLine(l).WithFilename(filename).AddErr(err)
			}

			continue
		}

		// multiline begin
		if l.isMultilineBegin() && !st.multiLine.isOn() {
			st.appendMultiLine(l)

			continue
		}

		// multiline end
		if st.multiLine.isOn() && l.isMultilineEnd(st.multiLine.isHeredoc()) {
			st.appendMultiLine(l)
			l = st.joinAnClearMultiLine()
		}

		// feed the multiline
		if st.multiLine.isOn() {
			st.appendMultiLine(l)

			continue
		}

		// section declaration
		s, ok := l.toSection()
		if ok {
			st.unsetCurrentLabel()
			st.setSection(s)

			continue
		}

		switch st.section {
		case sectionVars:
			name, value, ok := l.toVar()
			if ok {
				a.Vars.Set(name, value)
				st.setDef(name, sectionVars)

				continue
			}

			return ErrWrongVariableDeclaration.WithStack(st.stack).WithLine(l).WithFilename(filename)
		case sectionWords:
			w, ok := l.toWord()
			if !ok {
				return ErrWrongWordDeclaration.WithStack(st.stack).WithLine(l).WithFilename(filename)
			}

			_ = a.Vocabulary.Add(w.Label, w.Type, w.Synonyms...)
			st.setDef(w.Label, sectionWords)

			continue
		case sectionSysMsg:
			m, ok := l.toMsg(msg.SystemMsg)
			if !ok {
				return ErrWrongMessageDeclaration.WithStack(st.stack).WithLine(l).WithFilename(filename)
			}

			if err := a.Messages.Add(m); err != nil {
				return ErrWrongMessageDeclaration.WithStack(st.stack).AddErr(err).WithLine(l).WithFilename(filename)
			}
			st.setDef(m.Label, sectionSysMsg)

			continue
		case sectionUserMsgs:
			m, ok := l.toMsg(msg.UserMsg)
			if !ok {
				return ErrWrongMessageDeclaration.WithStack(st.stack).WithLine(l).WithFilename(filename)
			}

			if err := a.Messages.Add(m); err != nil {
				return ErrWrongMessageDeclaration.WithStack(st.stack).AddErr(err).WithLine(l).WithFilename(filename)
			}
			st.setDef(m.Label, sectionUserMsgs)

			continue
		case sectionItems:
			if st.hasCurrentLabel() {
				i := a.Objects.Get(st.currentLabel)
				if i == nil {
					return ErrWrongItemDeclaration.WithStack(st.stack).WithLine(l).WithFilename(filename)
				}

				desc, ok := l.labelAndTextRg("desc")
				if ok {
					i.SetDescription(desc)

					continue
				}

				desc, ok = l.labelAndTextRg("description")
				if ok {
					i.SetDescription(desc)

					continue
				}

				if l.optimized() == "is wearable" {
					i.SetWearable()
					continue
				}

				if l.optimized() == "is worn" {
					i.SetWearable()
					i.Wear()
					continue
				}

				if l.optimized() == "is created" {
					i.Create()
					continue
				}

				if l.optimized() == "is container" {
					i.SetContainer()
					continue
				}

				if l.optimized() == "is held" {
					i.Hold()
					continue
				}

				if itemLocationRg.MatchString(l.optimized()) {
					m := itemLocationRg.FindStringSubmatch(l.optimized())

					inLoc := a.Locations.Get(m[1])
					if inLoc == nil {
						inLoc = a.Locations.Set(m[1], loc.Undefined, loc.Undefined)
						st.setUndef(m[1], sectionLocs, l, filename)
					}

					i.SetLocation(inLoc)
					continue
				}

				if itemWeightRg.MatchString(l.optimized()) {
					m := itemWeightRg.FindStringSubmatch(l.optimized())
					w, err := strconv.Atoi(m[1])
					if err != nil {
						return ErrWrongItemWeight.WithStack(st.stack).WithLine(l).WithFilename(filename)
					}

					i.SetWeight(w)
					continue
				}

				if itemMaxWeightRg.MatchString(l.optimized()) {
					m := itemMaxWeightRg.FindStringSubmatch(l.optimized())
					w, err := strconv.Atoi(m[1])
					if err != nil {
						return ErrWrongItemWeight.WithStack(st.stack).WithLine(l).WithFilename(filename)
					}

					i.SetMaxWeight(w)
					continue
				}
			}

			label, noun, adj, ok := l.toItemDeclaration()
			if ok {
				if a.Objects.Exists(label) {
					return ErrDuplicatedItemLabel.WithStack(st.stack).WithLine(l).WithFilename(filename)
				}

				st.setCurrentLabel(label)
				st.setDef(label, sectionItems)

				nounWord := a.Vocabulary.Get(voc.Noun, noun)
				if nounWord == nil {
					nounWord = a.Vocabulary.Add(noun, voc.Noun)
					st.setUndef(noun, sectionWords, l, filename)
				}

				adjWord := a.Vocabulary.Get(voc.Adjective, adj)
				if adjWord == nil {
					adjWord = a.Vocabulary.Add(adj, voc.Adjective)
					st.setUndef(adj, sectionWords, l, filename)
				}

				if err := a.Objects.Add(obj.New(label, nounWord, adjWord)); err != nil {
					return ErrWrongItemDeclaration.WithStack(st.stack).AddErr(err).WithLine(l).WithFilename(filename)
				}

				st.setDef(label, sectionItems)

				continue
			}

			return ErrWrongItemDeclaration.WithStack(st.stack).WithLine(l).WithFilename(filename)
		case sectionLocs:
			label, ok := l.toLocationLabel()
			if ok {
				st.setCurrentLabel(label)
				a.Locations.Set(label, loc.Undefined, loc.Undefined)
				st.setDef(label, sectionLocs)

				continue
			} else {
				if !st.hasCurrentLabel() {
					return ErrWrongLocationLabelDeclaration.WithStack(st.stack).WithLine(l).WithFilename(filename)
				}
			}

			currentLocation := a.Locations.Get(st.currentLabel)

			desc, ok := l.toLocationDescription()
			if ok {
				currentLocation.Description = desc

				continue
			}

			title, ok := l.toLocationTitle()
			if ok {
				currentLocation.Title = title

				continue
			}

			exitMap, ok := l.toLocationConns()
			if ok {
				for wordLabel, destLabel := range exitMap {
					word := a.Vocabulary.FirstWithTypes(wordLabel, voc.Verb, voc.Noun)
					if word == nil {
						word = a.Vocabulary.Add(wordLabel, voc.Unknown)
						st.setUndef(wordLabel, sectionWords, l, filename)
					}

					dest := a.Locations.Get(destLabel)
					if dest == nil {
						dest = a.Locations.Set(destLabel, loc.Undefined, loc.Undefined)
						st.setUndef(destLabel, sectionLocs, l, filename)
					}

					currentLocation.SetConn(word, dest)
				}

				continue
			}
		case sectionProcs:
			// TODO
		default:
			return ErrOutOfSection.WithStack(st.stack).WithLine(l).WithFilename(filename)
		}

		// unmatched line
		return ErrUnknownDeclaration.WithStack(st.stack).WithLine(l).WithFilename(filename)
	}

	// check if there is an unclosed comment
	if st.comment.isOn() {
		return ErrUnclosedComment.WithStack(st.stack).WithLine(st.comment.lines[0]).WithFilename(filename)
	}

	// unclosed multiline
	if st.multiLine.isOn() {
		return ErrUnclosedMultiline.WithStack(st.stack).WithLine(st.multiLine.lines[0]).WithFilename(filename)
	}

	return nil
}
