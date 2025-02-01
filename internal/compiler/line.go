package compiler

import (
	"regexp"
	"strings"

	"thenewquill/internal/adventure/msg"
	"thenewquill/internal/adventure/voc"
)

const (
	labelMatcher = `[0-9\p{L}\-_]+`
	wordMatcher  = `[\p{L}\-_]+`
)

type line struct {
	text string
	n    int
}

func newLine(text string, n int) line {
	return line{text: text, n: n}
}

func (l line) optimized() string {
	spaceRg := regexp.MustCompile(`\s+`)
	commentRg := regexp.MustCompile(`//.*|/\*.*\*\/`)

	return strings.TrimSpace(spaceRg.ReplaceAllString(commentRg.ReplaceAllString(l.text, ""), " "))
}

func (l line) isCommentBegin() bool {
	re := regexp.MustCompile(`\s*/\*`)

	return re.MatchString(l.optimized())
}

func (l line) isCommentEnd() bool {
	re := regexp.MustCompile(`\*/\s*$`)

	return re.MatchString(l.optimized())
}

func (l line) isBlank() bool {
	return l.optimized() == ""
}

func (l line) toInClude() (string, bool) {
	re := regexp.MustCompile(`^INCLUDE\s+"(.*)"$`)

	if !re.MatchString(l.optimized()) {
		return "", false
	}

	return re.FindStringSubmatch(l.optimized())[1], true
}

func (l line) toSection() (section, bool) {
	re := regexp.MustCompile(`^SECTION\s+([\p{L}\s]+)$`)

	if !re.MatchString(l.optimized()) {
		return sectionNone, false
	}

	return sectionFromString(re.FindStringSubmatch(l.optimized())[1]), true
}

func (l line) toVar() (string, string, bool) {
	re := regexp.MustCompile(`^(\w+)\s*=\s*(.*)$`)

	o := l.optimized()

	if !re.MatchString(o) {
		return "", "", false
	}

	return re.FindStringSubmatch(o)[1], re.FindStringSubmatch(o)[2], true
}

func (l line) toWord() (voc.Word, bool) {
	re := regexp.MustCompile(`^(` + labelMatcher + `):(\s*((` + wordMatcher + `)),?)*$`)
	o := l.optimized()

	if !re.MatchString(o) {
		return voc.Word{}, false
	}

	w := voc.Word{}

	parts := strings.Split(o, ":")
	if len(parts) != 2 {
		return voc.Word{}, false
	}

	w.Type = voc.WordTypeFromString(parts[0])

	words := strings.Split(parts[1], ",")
	for i, word := range words {
		if i == 0 {
			w.Label = strings.TrimSpace(word)

			continue
		}

		w.Synonyms = append(w.Synonyms, strings.TrimSpace(word))
	}

	return w, true
}

func (l line) toMsg(t msg.MsgType) (msg.Msg, bool) {
	re := regexp.MustCompile(`^([0-9\p{L}-_]+):\s+["^(\\")]{1}(.+)["^(\\")]{1}$`)
	o := l.optimized()

	if !re.MatchString(o) {
		return msg.Msg{}, false
	}

	parts := re.FindStringSubmatch(o)
	if len(parts) != 3 {
		return msg.Msg{}, false
	}

	return msg.Msg{Type: t, Label: parts[1], Text: parts[2]}, true
}
