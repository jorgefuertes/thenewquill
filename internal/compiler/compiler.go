package compiler

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/processor"
	"github.com/jorgefuertes/thenewquill/internal/compiler/section"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
)

const VERSION = "1.0.0"

func Compile(filename string) (*adventure.Adventure, error) {
	a := adventure.New()
	st := status.New()

	err := compileFile(st, filename, a)
	cErr, ok := err.(cerr.CompilerError)
	if ok {
		fmt.Println(cErr.Dump())
		return a, cErr
	}

	// check for unresolved labels
	for _, udf := range st.Undefs {
		fmt.Println(
			cerr.ErrUnresolvedLabel.WithFilename(udf.File).
				WithLine(udf.Line).
				WithSection(udf.Section).
				AddMsgf("%s `%s` remains undefined", udf.Section.String(), udf.Label).
				Dump(),
		)
		err = cerr.ErrRemainingUnresolvedLabels
	}
	if err != nil {
		return a, err
	}

	return a, a.Validate()
}

func compileFile(st *status.Status, filename string, a *adventure.Adventure) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	st.PushFilename(filename)
	defer st.PopFilename()

	n := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		n++
		l := line.New(scanner.Text(), n)

		if l.IsBlank() {
			continue
		}

		st.AppendStack(l)

		if l.IsOneLineComment() {
			continue
		}

		// comments are ignored
		if l.IsCommentBegin() && !st.Comment.IsOn() {
			st.SetComment(l)

			continue
		}

		if st.Comment.IsOn() && l.IsCommentEnd() {
			st.UnsetComment()

			continue
		}

		if st.Comment.IsOn() {
			continue
		}

		// follow includes
		f, ok := l.AsInclude()
		if ok {
			err := compileFile(st, filepath.Dir(st.CurrentFilename())+"/"+f, a)
			if err != nil {
				cErr, ok := err.(cerr.CompilerError)
				if ok {
					return cErr
				}

				return cerr.ErrCannotOpenIncludedFile.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
					WithFilename(st.CurrentFilename()).AddErr(err)
			}

			continue
		}

		// multiline begin
		if l.IsMultilineBegin() && !st.MultiLine.IsOn() {
			st.AppendLine(l)

			continue
		}

		// multiline end
		if st.MultiLine.IsOn() && l.IsMultilineEnd(st.MultiLine.IsHeredoc()) {
			st.AppendLine(l)
			l = st.MultiLine.Join()
			st.MultiLine.Clear()
		}

		// feed the multiline
		if st.MultiLine.IsOn() {
			st.AppendLine(l)

			continue
		}

		// section declaration
		s, ok := l.AsSection()
		if ok {
			st.CurrentLabel = ""
			st.Section = s

			continue
		}

		if st.Section == section.None {
			return cerr.ErrOutOfSection.WithStack(st.Stack).WithSection(st.Section).WithLine(l).
				WithFilename(st.CurrentFilename())
		}

		err := processor.ProcessLine(l, st, a)
		if err != nil {
			return err
		} else {
			continue
		}
	}

	// check if there is an unclosed comment
	if st.Comment.IsOn() {
		l, _ := st.Comment.GetByIndex(0)
		return cerr.ErrUnclosedComment.WithSection(st.Section).WithStack(st.Stack).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	// unclosed multiline
	if st.MultiLine.IsOn() {
		l, _ := st.MultiLine.GetByIndex(0)
		return cerr.ErrUnclosedMultiline.WithSection(st.Section).WithStack(st.Stack).WithLine(l).
			WithFilename(st.CurrentFilename())
	}

	return nil
}
