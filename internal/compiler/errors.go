package compiler

var (
	ErrOutOfSection                  = newCompilerError("line outside of section")
	ErrUnclosedComment               = newCompilerError("unclosed comment")
	ErrUnclosedString                = newCompilerError("unclosed string")
	ErrCannotOpenIncludedFile        = newCompilerError("cannot open included file")
	ErrUnknownDeclaration            = newCompilerError("unknown declaration")
	ErrWrongVariableDeclaration      = newCompilerError("wrong variable declaration")
	ErrWrongWordDeclaration          = newCompilerError("wrong word declaration")
	ErrCannotCreateWord              = newCompilerError("cannot create word")
	ErrWrongLocationLabelDeclaration = newCompilerError("wrong location label declaration")
	ErrWrongMessageDeclaration       = newCompilerError("wrong message declaration")
	ErrWrongExitsDeclaration         = newCompilerError("wrong location exits declaration")
	ErrUnclosedMultiline             = newCompilerError("unclosed multiline declaration")
	ErrUnresolvedLabel               = newCompilerError("unresolved label")
	ErrRemainingUnresolvedLabels     = newCompilerError("remaining unresolved labels")
)
