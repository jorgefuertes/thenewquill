package compiler_error

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
)

type CompilerError struct {
	section  kind.Kind
	filename string
	stack    []line.Line
	l        line.Line
	msgs     []string
}

func (e CompilerError) Dump() {
	o := NewOutput("COMPILATION ERROR")

	if e.Is(ErrValidation) {
		o = NewOutput("VALIDATION ERROR")
		for _, m := range e.msgs {
			if m == "validation error" {
				continue
			}

			o.addLine(-1, m)
			o.addNL()
		}

		o.Print()

		return
	}

	if len(e.stack) > 0 {
		for _, l := range e.stack {
			if l.Number() == e.l.Number() {
				break
			}

			o.addLine(l.Num, l.Text)
		}
	}

	if e.filename != "" {
		o.addLine(-1, downArrowHead+" FILE "+e.filename)
	}

	if e.section != kind.None {
		o.addLine(-1, downArrowHead+" SECTION "+e.section.TitleString())
	}

	o.addNL()
	o.addLine(e.l.Num, e.l.Text)
	o.addNL()

	for _, msg := range e.msgs {
		o.addLine(-1, upArrowHead+" ERROR "+msg)
	}

	o.Print()
}

func (e CompilerError) Error() string {
	return e.msgs[0]
}

func newCompilerError(msg string) CompilerError {
	return CompilerError{
		msgs:    []string{msg},
		section: kind.None,
	}
}

func (e CompilerError) WithSection(k kind.Kind) CompilerError {
	e.section = k

	return e
}

func (e CompilerError) WithLine(l line.Line) CompilerError {
	e.l = l

	return e
}

func (e CompilerError) WithStack(s []line.Line) CompilerError {
	e.stack = s

	if len(s) > 0 {
		e.l = s[len(s)-1]
	}

	return e
}

func (e CompilerError) WithFilename(filename string) CompilerError {
	e.filename = filename

	return e
}

func (e CompilerError) AddMsg(msg string) CompilerError {
	e.msgs = append(e.msgs, msg)

	return e
}

func (e CompilerError) AddMsgf(format string, args ...any) CompilerError {
	e.msgs = append(e.msgs, fmt.Sprintf(format, args...))

	return e
}

func (e CompilerError) AddErr(err error) CompilerError {
	e.msgs = append(e.msgs, err.Error())

	return e
}

func (e CompilerError) Is(err error) bool {
	for _, msg := range e.msgs {
		if msg == err.Error() {
			return true
		}
	}

	return false
}

func (e CompilerError) IsOK() bool {
	return e.Is(OK)
}
