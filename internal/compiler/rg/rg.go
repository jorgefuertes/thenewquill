package rg

import (
	"regexp"
)

var (
	Label          = regexp.MustCompile(`^([\d\p{L}\-_]+)$`)
	InlineComment  = regexp.MustCompile(`(//.*|/\*.*)$`)
	OneLinecomment = regexp.MustCompile(`^\s*(/\*.*\*/|//.*)\s*$`)
	CommentBegin   = regexp.MustCompile(`^\s*/\*`)
	CommentEnd     = regexp.MustCompile(`\*/\s*$`)
	MultilineBegin = regexp.MustCompile(`("""\s*)$`)
	MultilineEnd   = regexp.MustCompile(`^\s*"""\s*$`)
	Indent         = regexp.MustCompile(`^(\s*)`)
	Continue       = regexp.MustCompile(`(\\\s*)$`)

	Blank   = regexp.MustCompile(`^\s*$`)
	Include = regexp.MustCompile(`^INCLUDE\s+"(.*)"$`)
	Section = regexp.MustCompile(`(?i)^SECTION\s+([\p{L}\s]+)$`)
	Var     = regexp.MustCompile(`^([\d\p{L}\-_]+)\s*=\s*"?((?:[^"]|.\\")+)"?$`)
	Float   = regexp.MustCompile(`^(\d+\.\d+)$`)
	Int     = regexp.MustCompile(`^(\d+)$`)
	Bool    = regexp.MustCompile(`^(true|false)$`)
	Word    = regexp.MustCompile(`^([\d\p{L}\-_]+):\s*((?:[\d\p{L}\-_^]+),*\s*)+$`)
	Msg     = regexp.MustCompile(
		`(?s)^((?:[\d\p{L}\-_]+)(?:\.(?:zero|one|many)){0,1}):\s+["^(\\")]{1}(.+)["^(\\")]{1}$`,
	)
	MsgPlural = regexp.MustCompile(`(?s)^([\d\p{L}\-_]+)\.(zero|one|many):\s+["^(\\")]{1}(.+)["^(\\")]{1}$`)
	LocLabel  = regexp.MustCompile(`^\s*([\d\p{L}\-_]+):\s*$`)
	LocConns  = regexp.MustCompile(
		`^\s*(exits|conns|connections)\s*:\s*(\s*([\d\p{L}\-_]+\s+[\d\p{L}\-_]+\s*,?))+.?$`,
	)
	LabelNounAdjDeclaration = regexp.MustCompile(`^\s*([\d\p{L}\-_]+):\s+([\d\p{L}\-_]+)\s+([\d\p{L}\-_]+)\s*$`)
)

func IsValidLabel(label string) bool {
	return Label.MatchString(label)
}
