package line

import (
	"regexp"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/compiler/rg"
)

type Line struct {
	text string
	n    int
}

func New(text string, n int) Line {
	return Line{text: text, n: n}
}

func (l Line) OptimizedText() string {
	return strings.TrimSpace(rg.InlineComment.ReplaceAllString(l.text, ""))
}

func (l Line) Text() string {
	return l.text
}

func (l *Line) Add(text string) {
	l.text += text
}

func (l Line) Number() int {
	return l.n
}

// GetTextForLabelName returns the text for the given label and true if it was found
func (l Line) GetTextForLabel(label string) (string, bool) {
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

func (l Line) GetTextForFirstFoundLabel(labels ...string) (string, bool) {
	for _, label := range labels {
		text, ok := l.GetTextForLabel(label)
		if ok {
			return text, ok
		}
	}

	return "", false
}
