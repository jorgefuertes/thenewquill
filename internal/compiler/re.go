package compiler

import (
	"regexp"
	"strings"
)

const (
	labelMatcher = `[0-9\p{L}\-_]+`
	wordMatcher  = `[\p{L}\-_]+`
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
