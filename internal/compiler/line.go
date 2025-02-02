package compiler

import (
	"regexp"
	"strconv"
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
	commentRg := regexp.MustCompile(`(//.*|/\*.*)$`)

	return strings.TrimSpace(commentRg.ReplaceAllString(l.text, ""))
}

func (l line) isCommentBegin() bool {
	re := regexp.MustCompile(`^\s*/\*`)

	return re.MatchString(l.text)
}

func (l line) isCommentEnd() bool {
	re := regexp.MustCompile(`\*/\s*$`)

	return re.MatchString(l.optimized())
}

func (l line) isBlank() bool {
	blankRg := regexp.MustCompile(`^\s*$`)

	return blankRg.MatchString(l.text)
}

func (l line) isOneLineComment() bool {
	commentRg := regexp.MustCompile(`^\s*(/\*.*\*/|//.*)\s*$`)

	return commentRg.MatchString(l.text)
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

func (l line) toVar() (string, any, bool) {
	re := regexp.MustCompile(`^(\w+)\s*=\s*"?((?:[^"]|.\\")+)"?$`)

	o := l.optimized()

	if !re.MatchString(o) {
		return "", "", false
	}

	name := re.FindStringSubmatch(o)[1]
	valueStr := re.FindStringSubmatch(o)[2]

	floatRg := regexp.MustCompile(`^(\d+\.\d+)$`)
	intRg := regexp.MustCompile(`^(\d+)$`)
	boolRg := regexp.MustCompile(`^(true|false)$`)

	if floatRg.MatchString(valueStr) {
		value, _ := strconv.ParseFloat(valueStr, 64)

		return name, value, true
	}

	if intRg.MatchString(valueStr) {
		value, _ := strconv.ParseInt(valueStr, 10, 64)

		return name, value, true
	}

	if boolRg.MatchString(valueStr) {
		value, _ := strconv.ParseBool(valueStr)

		return name, value, true
	}

	return name, valueStr, true
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
	re := regexp.MustCompile(`(?s)^([0-9\p{L}-_]+):\s+["^(\\")]{1}(.+)["^(\\")]{1}$`)

	if !re.MatchString(l.text) {
		return msg.Msg{}, false
	}

	parts := re.FindStringSubmatch(l.text)
	if len(parts) != 3 {
		return msg.Msg{}, false
	}

	return msg.Msg{Type: t, Label: parts[1], Text: parts[2]}, true
}

func (l line) getIndent() string {
	re := regexp.MustCompile(`^(\s+)`)

	if !re.MatchString(l.text) {
		return ""
	}

	return re.FindStringSubmatch(l.text)[1]
}

func (l line) isMultilineBegin() bool {
	return multilineBeginRg.MatchString(l.text)
}

func (l line) isMultilineEnd() bool {
	return multilineEndRg.MatchString(l.text)
}
