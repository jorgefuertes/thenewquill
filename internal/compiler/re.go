package compiler

import "regexp"

const (
	labelMatcher = `[0-9\p{L}\-_]+`
	wordMatcher  = `[\p{L}\-_]+`
)

func (l line) labelAndTextRg(label string) (string, bool) {
	re := regexp.MustCompile(`(?s)^\s*` + label + `:\s+["^(\\")]{1}(.+)["^(\\")]{1}`)

	if !re.MatchString(l.text) {
		return "", false
	}

	return re.FindStringSubmatch(l.text)[1], true
}
