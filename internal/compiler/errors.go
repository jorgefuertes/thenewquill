package compiler

import "fmt"

type compilerError struct {
	filename string
	l        line
	msgs     []string
}

func (e compilerError) Error() string {
	output := fmt.Sprintf("[ERROR] 🔻 compiling file: '%s'\n[%05d] %s", e.filename, e.l.n, e.l.text)
	for _, msg := range e.msgs {
		output += fmt.Sprintf("\n[ERROR] ↪ %s", msg)
	}

	return output
}

func newCompilerError(msg string) compilerError {
	return compilerError{msgs: []string{msg}}
}

func (e compilerError) WithLine(l line) compilerError {
	e.l = l
	return e
}

func (e compilerError) WithFilename(filename string) compilerError {
	e.filename = filename
	return e
}

func (e compilerError) AddMsg(msg string) compilerError {
	e.msgs = append(e.msgs, msg)

	return e
}

func (e compilerError) AddErr(err error) compilerError {
	e.msgs = append(e.msgs, err.Error())

	return e
}

var (
	ErrOutOfSection             = newCompilerError("line outside of section")
	ErrUnclosedComment          = newCompilerError("unclosed comment")
	ErrUnclosedString           = newCompilerError("unclosed string")
	ErrCannotOpenIncludedFile   = newCompilerError("cannot open included file")
	ErrUnknownDeclaration       = newCompilerError("unknown declaration")
	ErrWrongVariableDeclaration = newCompilerError("wrong variable declaration")
	ErrWrongWordDeclaration     = newCompilerError("wrong word declaration")
	ErrWrongMessageDeclaration  = newCompilerError("wrong message declaration")
	ErrUnclosedMultiline        = newCompilerError("unclosed multiline declaration")
)
