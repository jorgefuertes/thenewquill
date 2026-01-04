package line

import (
	"regexp"
	"strings"

	"github.com/jorgefuertes/thenewquill/internal/compiler/rg"
)

type Multi struct {
	lines []Line
}

func NewMulti(lines ...Line) Multi {
	return Multi{lines: lines}
}

func (m *Multi) Clear() {
	m.lines = []Line{}
}

// Append appends a line
func (m *Multi) Append(l Line) {
	m.lines = append(m.lines, l)
}

// IsOn returns true if there are lines
func (m Multi) IsOn() bool {
	return len(m.lines) > 0
}

func (m Multi) GetByIndex(i int) (Line, bool) {
	if i < 0 || i >= len(m.lines) {
		return Line{}, false
	}

	return m.lines[i], true
}

func (m Multi) Len() int {
	return len(m.lines)
}

// IsHeredoc returns true if the first line is a heredoc or no lines are present
func (m Multi) IsHeredoc() bool {
	if len(m.lines) == 0 {
		return true
	}

	return rg.MultilineBegin.MatchString(m.lines[0].Text)
}

// GetIndent returns the indentation text of the first line
func (m Multi) GetIndent() string {
	if len(m.lines) < 2 {
		return ""
	}

	return rg.Indent.FindStringSubmatch(m.lines[1].Text)[1]
}

func (m Multi) Join() Line {
	var output string

	for i, l := range m.lines {
		current := l.Text
		cont := false

		if i == 0 {
			current = rg.MultilineBegin.ReplaceAllString(current, `"`)
			cont = true
		} else {
			current = strings.Replace(current, m.GetIndent(), "", 1)
		}

		if rg.MultilineEnd.MatchString(current) {
			current = `"`
			output = strings.TrimSuffix(output, "\n")
		}

		if rg.Continue.MatchString(current) {
			cont = true
			current = rg.Continue.ReplaceAllString(current, "")
			current = regexp.MustCompile(`^(\s*)`).ReplaceAllString(current, "")
		}

		output += current

		if !cont && i != m.Len()-1 && m.lines[i+1].Text != `"""` {
			output += "\n"
		}
	}

	return Line{Text: output, Num: m.lines[0].Number()}
}
