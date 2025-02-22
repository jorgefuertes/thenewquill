package compiler

import (
	"strconv"
	"strings"

	"thenewquill/internal/adventure/msg"
	"thenewquill/internal/adventure/voc"
)

type line struct {
	text string
	n    int
}

func newLine(text string, n int) line {
	return line{text: text, n: n}
}

func (l line) optimized() string {
	return strings.TrimSpace(inlineCommentRg.ReplaceAllString(l.text, ""))
}

func (l line) isCommentBegin() bool {
	return commentBeginRg.MatchString(l.text)
}

func (l line) isCommentEnd() bool {
	return commentEndRg.MatchString(l.optimized())
}

func (l line) isBlank() bool {
	return blankRg.MatchString(l.text)
}

func (l line) isOneLineComment() bool {
	return oneLinecommentRg.MatchString(l.text)
}

func (l line) toInClude() (string, bool) {
	if !includeRg.MatchString(l.optimized()) {
		return "", false
	}

	return includeRg.FindStringSubmatch(l.optimized())[1], true
}

func (l line) toSection() (section, bool) {
	if !sectionRg.MatchString(l.optimized()) {
		return sectionNone, false
	}

	return sectionFromString(sectionRg.FindStringSubmatch(l.optimized())[1]), true
}

func (l line) toVar() (string, any, bool) {
	o := l.optimized()

	if !varRg.MatchString(o) {
		return "", "", false
	}

	name := varRg.FindStringSubmatch(o)[1]
	valueStr := varRg.FindStringSubmatch(o)[2]

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
	o := l.optimized()

	if !wordRg.MatchString(o) {
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

func (l line) toMsg(t msg.MsgType) (*msg.Msg, bool) {
	if !msgRg.MatchString(l.text) {
		return nil, false
	}

	parts := msgRg.FindStringSubmatch(l.text)
	if len(parts) != 3 {
		return nil, false
	}

	return msg.New(t, parts[1], parts[2]), true
}

func (l line) isMultilineBegin() bool {
	return multilineBeginRg.MatchString(l.text) || continueRg.MatchString(l.text)
}

func (l line) isMultilineEnd(isHeredoc bool) bool {
	if isHeredoc {
		return multilineEndRg.MatchString(l.text)
	}

	return !continueRg.MatchString(l.text)
}

func (l line) toLocationLabel() (string, bool) {
	if !locLabelRg.MatchString(l.text) {
		return "", false
	}

	return locLabelRg.FindStringSubmatch(l.text)[1], true
}

func (l line) toLocationDescription() (string, bool) {
	return l.labelAndTextRg("desc")
}

func (l line) toLocationTitle() (string, bool) {
	return l.labelAndTextRg("title")
}

func (l line) toLocationConns() (map[string]string, bool) {
	exits := make(map[string]string, 0)

	if !locConnsRg.MatchString(l.text) {
		return exits, false
	}

	parts := strings.Split(strings.Split(l.optimized(), ":")[1], ",")

	for _, part := range parts {
		words := strings.Split(strings.TrimSpace(part), " ")
		if len(words) != 2 {
			return exits, false
		}

		exits[words[0]] = words[1]
	}

	return exits, true
}

func (l line) toItemDeclaration() (label, noun, adjetive string, ok bool) {
	if !itemDeclarationRg.MatchString(l.text) {
		return "", "", "", false
	}

	m := itemDeclarationRg.FindStringSubmatch(l.text)
	label = m[1]
	noun = m[2]
	adjetive = m[3]

	return label, noun, adjetive, true
}
