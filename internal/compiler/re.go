package compiler

import (
	"regexp"
	"strings"
)

var (
	inlineCommentRg  = regexp.MustCompile(`(//.*|/\*.*)$`)
	oneLinecommentRg = regexp.MustCompile(`^\s*(/\*.*\*/|//.*)\s*$`)
	commentBeginRg   = regexp.MustCompile(`^\s*/\*`)
	commentEndRg     = regexp.MustCompile(`\*/\s*$`)
	blankRg          = regexp.MustCompile(`^\s*$`)
	includeRg        = regexp.MustCompile(`^INCLUDE\s+"(.*)"$`)
	sectionRg        = regexp.MustCompile(`(?i)^SECTION\s+([\p{L}\s]+)$`)
	varRg            = regexp.MustCompile(`^([\d\p{L}\-_]+)\s*=\s*"?((?:[^"]|.\\")+)"?$`)
	floatRg          = regexp.MustCompile(`^(\d+\.\d+)$`)
	intRg            = regexp.MustCompile(`^(\d+)$`)
	boolRg           = regexp.MustCompile(`^(true|false)$`)
	wordRg           = regexp.MustCompile(`^([\d\p{L}\-_]+):\s*((?:[\d\p{L}\-_^]+),*\s*)+$`)
	msgRg            = regexp.MustCompile(
		`(?s)^((?:[\d\p{L}\-_]+)(?:\.(?:zero|one|more)){0,1}):\s+["^(\\")]{1}(.+)["^(\\")]{1}$`,
	)
	msgPluralRg = regexp.MustCompile(`(?s)^([\d\p{L}\-_]+)\.(zero|one|more):\s+["^(\\")]{1}(.+)["^(\\")]{1}$`)
	locLabelRg  = regexp.MustCompile(`^\s*([\d\p{L}\-_]+):\s*$`)
	locConnsRg  = regexp.MustCompile(
		`^\s*(exits|conns|connections)\s*:\s*(\s*([\d\p{L}\-_]+\s+[\d\p{L}\-_]+\s*,?))+.?$`,
	)
	itemDeclarationRg = regexp.MustCompile(`^\s*([\d\p{L}\-_]+):\s+([\d\p{L}\-_]+)\s+([\d\p{L}\-_]+)\s*$`)
	itemLocationRg    = regexp.MustCompile(`^\s*is\s+in\s+([\d\p{L}\-_]+)\s*$`)
	itemWeightRg      = regexp.MustCompile(`^\s*has\s+weight\s+(\d+)\s*$`)
	itemMaxWeightRg   = regexp.MustCompile(`^\s*has\s+max\s+weight\s+(\d+)\s*$`)
)

func (l line) labelAndTextRg(label string) (string, bool) {
	re := regexp.MustCompile(`(?s)^\s*` + label + `:\s+["^(\\")]{1}(.+)["^(\\")]{1}`)

	if !re.MatchString(l.text) {
		return "", false
	}

	text := re.FindStringSubmatch(l.text)[1]

	// normalize escaped quotes
	text = strings.ReplaceAll(text, `\"`, `"`)
	text = strings.ReplaceAll(text, `\'`, `'`)

	return text, true
}
