package compiler

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"

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
		if l.isCommentBegin() && !st.inComment {
			st.setComment(l)

			continue
		}

		if l.isCommentEnd() && st.inComment {
			st.unsetComment()

			continue
		}

		if st.inComment {
			continue
		}

		if l.isBlank() {
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

				return ErrCannotOpenIncludedFile.WithLine(l).WithFilename(filename).AddErr(err)
			}

			continue
		}

		// multiline end
		if st.inMulti {
			multilineEndRg := regexp.MustCompile(`^(\s*)"""`)
			if multilineEndRg.MatchString(l.text) {
				l = st.getMultiLine()
				l.text += `"`
				l.n = n
				st.unsetMultiLine()
			} else {
				st.appendMultiLine(l)

				continue
			}
		}

		// mline begin
		multilineBeginRg := regexp.MustCompile(`("""(\s*))$`)
		if multilineBeginRg.MatchString(l.text) {
			l.text = multilineBeginRg.ReplaceAllString(l.text, `"`)
			st.setMultiLine(l)

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

			return ErrWrongVariableDeclaration.WithLine(l).WithFilename(filename)
		case sectionWords:
			w, ok := l.toWord()
			if !ok {
				return ErrWrongWordDeclaration.WithLine(l).WithFilename(filename)
			}

			err := a.Vocabulary.Add(w.Label, w.Type, w.Synonyms...)
			if err != nil {
				return ErrWrongWordDeclaration.AddErr(err).WithLine(l).WithFilename(filename)
			}

			continue
		case sectionSysMsg:
			m, ok := l.toMsg(msg.SystemMsg)
			if !ok {
				return ErrWrongMessageDeclaration.WithLine(l).WithFilename(filename)
			}

			if err := a.Messages.Add(m); err != nil {
				return ErrWrongMessageDeclaration.AddErr(err).WithLine(l).WithFilename(filename)
			}

			continue
		case sectionUserMsgs:
			m, ok := l.toMsg(msg.UserMsg)
			if !ok {
				return ErrWrongMessageDeclaration.WithLine(l).WithFilename(filename)
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
			return ErrOutOfSection.WithLine(l).WithFilename(filename)
		}

		// unmatched line
		return ErrUnknownDeclaration.WithLine(l).WithFilename(filename)
	}

	// check if there is an unclosed comment
	if st.inComment {
		return ErrUnclosedComment.WithLine(st.getLastLine()).WithFilename(filename)
	}

	// unclosed multiline
	if st.inMulti {
		return ErrUnclosedMultiline.WithLine(st.multiLine).WithFilename(filename)
	}

	return nil
}
