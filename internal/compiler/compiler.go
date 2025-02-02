package compiler

import (
	"bufio"
	"os"
	"path/filepath"

	"thenewquill/internal/adventure"
	"thenewquill/internal/adventure/msg"
)

func Compile(filename string) (*adventure.Adventure, error) {
	a := adventure.New()
	st := newStatus()

	err := compileFile(st, filename, a)

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

		// comments and blank lines are ignored
		if l.isCommentBegin() && !st.comment.isOn() {
			st.setComment(l)
			st.appendStack(l)

			continue
		}

		if l.isCommentEnd() && st.comment.isOn() {
			st.unsetComment()
			st.appendStack(l)

			continue
		}

		if st.comment.isOn() {
			continue
		}

		if l.isBlank() {
			continue
		}

		st.appendStack(l)

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
			st.startMultiLine(l)

			continue
		}

		// multiline end
		if l.isMultilineEnd() && st.multiLine.isOn() {
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
			st.setSection(s, l)

			continue
		}

		switch st.section {
		case sectionVars:
			name, value, ok := l.toVar()
			if ok {
				a.Vars.Set(name, value)

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

			continue
		case sectionSysMsg:
			m, ok := l.toMsg(msg.SystemMsg)
			if !ok {
				return ErrWrongMessageDeclaration.WithStack(st.stack).WithLine(l).WithFilename(filename)
			}

			if err := a.Messages.Add(m); err != nil {
				return ErrWrongMessageDeclaration.WithStack(st.stack).AddErr(err).WithLine(l).WithFilename(filename)
			}

			continue
		case sectionUserMsgs:
			m, ok := l.toMsg(msg.UserMsg)
			if !ok {
				return ErrWrongMessageDeclaration.WithStack(st.stack).WithLine(l).WithFilename(filename)
			}

			a.Messages.Add(m)

			continue
		case sectionObjs:
			// TODO
		case sectionLocs:
			// TODO
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
