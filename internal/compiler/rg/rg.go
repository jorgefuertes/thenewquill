package rg

import (
	"regexp"
)

const (
	labelGroup      = `([\d\p{L}\-_\.]+)`
	intOrFloatGroup = `(\d+|\d+\.\d+)`
)

var (
	Label          = regexp.MustCompile(`^` + labelGroup + `$`)
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
	Var     = regexp.MustCompile(`^` + labelGroup + `\s*=\s*"?((?:[^"]|.\\")+)"?$`)
	Float   = regexp.MustCompile(`^(\d+\.\d+)$`)
	Int     = regexp.MustCompile(`^(\d+)$`)
	Bool    = regexp.MustCompile(`^(true|false)$`)
	Word    = regexp.MustCompile(`^` + labelGroup + `:\s*((?:[\d\p{L}\-_^]+),*\s*)+$`)
	Msg     = regexp.MustCompile(
		`(?s)^((?:[\d\p{L}\-_]+)(?:\.(?:zero|one|many)){0,1}):\s+["^(\\")]{1}(.+)["^(\\")]{1}$`,
	)
	MsgPlural = regexp.MustCompile(`(?s)^` + labelGroup + `\.(zero|one|many):\s+["^(\\")]{1}(.+)["^(\\")]{1}$`)
	LocLabel  = regexp.MustCompile(`^\s*` + labelGroup + `:\s*$`)
	LocConns  = regexp.MustCompile(
		`^\s*(exits|conns|connections)\s*:\s*(\s*([\d\p{L}\-_]+\s+[\d\p{L}\-_]+\s*,?))+.?$`,
	)
	LabelNounAdjDeclaration = regexp.MustCompile(
		`^\s*` + labelGroup + `:\s+` + labelGroup + `\s+` + labelGroup + `\s*$`,
	)
	ItemAt        = regexp.MustCompile(`^\s*(is at|is in|is worn by)\s+` + labelGroup + `\s*$`)
	ItemWeight    = regexp.MustCompile(`^\s*(has weight|weight|weighs)\s+` + intOrFloatGroup + `\s*$`)
	ItemMaxWeight = regexp.MustCompile(`^\s*(has max weight|max weight)\s+` + intOrFloatGroup + `\s*$`)
	Blob          = regexp.MustCompile(`^\s*` + labelGroup + `:\s+(.+\/[^\/]+\.[^.]+)\s*$`)
)

func IsValidLabel(label string) bool {
	return Label.MatchString(label)
}
