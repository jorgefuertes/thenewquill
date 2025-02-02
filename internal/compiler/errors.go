package compiler

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
