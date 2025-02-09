package compiler

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"thenewquill/internal/adventure"
	"thenewquill/internal/adventure/msg"
)

func Compile(filename string) (*adventure.Adventure, error) {
	a := adventure.New()
	st := newStatus()

	err := compileFile(st, filename, a)
	cErr, ok := err.(compilerError)
	if ok {
		fmt.Println(cErr.Dump())
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
			st.setSection(s)

			continue
		}

		switch st.section {
		case sectionVars:
			name, value, ok := l.toVar()
			if ok {
				a.Vars.Set(name, value)
				st.appendLabel(name, sectionVars, true, l)

				continue
			}

			return ErrWrongVariableDeclaration.WithStack(st.stack).WithLine(l).WithFilename(filename)
		case sectionWords:
			w, ok := l.toWord()
			if !ok {
				return ErrWrongWordDeclaration.WithStack(st.stack).WithLine(l).WithFilename(filename)
			}

			err := a.Vocabulary.Add(w.Label, w.Type, w.Synonyms...)
			if err != nil {
				return ErrWrongWordDeclaration.WithStack(st.stack).AddErr(err).WithLine(l).WithFilename(filename)
			}

			st.appendLabel(w.Label, sectionWords, true, l)

			continue
		case sectionSysMsg:
			m, ok := l.toMsg(msg.SystemMsg)
			if !ok {
				return ErrWrongMessageDeclaration.WithStack(st.stack).WithLine(l).WithFilename(filename)
			}

			if err := a.Messages.Add(m); err != nil {
				return ErrWrongMessageDeclaration.WithStack(st.stack).AddErr(err).WithLine(l).WithFilename(filename)
			}

			st.appendLabel(m.Label, sectionSysMsg, true, l)

			continue
		case sectionUserMsgs:
			m, ok := l.toMsg(msg.UserMsg)
			if !ok {
				return ErrWrongMessageDeclaration.WithStack(st.stack).WithLine(l).WithFilename(filename)
			}

			if err := a.Messages.Add(m); err != nil {
				return ErrWrongMessageDeclaration.WithStack(st.stack).AddErr(err).WithLine(l).WithFilename(filename)
			}

			st.appendLabel(m.Label, sectionUserMsgs, true, l)

			continue
		case sectionObjs:
			// TODO
		case sectionLocs:
			if st.currentLocation.Label != "" {
				desc, ok := l.toLocationDescription()
				if ok {
					st.currentLocation.Description = desc

					continue
				}

				title, ok := l.toLocationTitle()
				if ok {
					st.currentLocation.Title = title

					continue
				}
			}

			label, ok := l.toLocationLabel()
			if ok {
				st.setCurrentLocation(label)
				st.appendLabel(label, sectionLocs, true, l)

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
