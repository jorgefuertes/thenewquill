package compiler

import (
	"fmt"

	"thenewquill/internal/compiler/line"
	"thenewquill/internal/compiler/section"
)

type CompilerError struct {
	section  section.Section
	filename string
	stack    []line.Line
	l        line.Line
	msgs     []string
}

func (e CompilerError) Dump() string {
	output := ""

	if len(e.stack) > 0 {
		for _, l := range e.stack {
			if l.Number() == e.l.Number() {
				break
			}

			output += fmt.Sprintf("[%05d] %s\n", l.Number(), l.Text())
		}
	}

	output += fmt.Sprintf(
		"[ERROR] 🔻 FILE \"%s\" ❗ SECTION \"%s\"\n[%05d] %s",
		e.filename,
		e.section.String(),
		e.l.Number(),
		e.l.Text(),
	)
	for _, msg := range e.msgs {
		output += fmt.Sprintf("\n[ERROR] 🔺 %s", msg)
	}

	return output
}

func (e CompilerError) Error() string {
	return e.msgs[0]
}

func newCompilerError(msg string) CompilerError {
	return CompilerError{msgs: []string{msg}}
}

func (e CompilerError) WithSection(s section.Section) CompilerError {
	e.section = s

	return e
}

func (e CompilerError) WithLine(l line.Line) CompilerError {
	e.l = l

	return e
}

func (e CompilerError) WithStack(s []line.Line) CompilerError {
	e.stack = s

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
