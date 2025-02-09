package compiler

import (
	"regexp"
	"strings"
)

var (
	multilineBeginRg = regexp.MustCompile(`("""\s*)$`)
	multilineEndRg   = regexp.MustCompile(`^\s*"""`)
	indentRg         = regexp.MustCompile(`^(\s*)`)
	continueRg       = regexp.MustCompile(`(\\\s*)$`)
)

type multi struct {
	lines []line
}

func (m *multi) append(l line) {
	m.lines = append(m.lines, l)
}

func (m multi) isOn() bool {
	return len(m.lines) > 0
}

func (m multi) isHeredoc() bool {
	if len(m.lines) == 0 {
		return true
	}

	return multilineBeginRg.MatchString(m.lines[0].text)
}

// returns true if its the end of the multiline
func (s *status) appendMultiLine(l line) {
	s.multiLine.append(l)
}

func (m multi) getIndent() string {
	if len(m.lines) < 2 {
		return ""
	}

	return indentRg.FindStringSubmatch(m.lines[1].text)[1]
}

func (s *status) joinAnClearMultiLine() line {
	var output string
	n := s.multiLine.lines[0].n

	for i, l := range s.multiLine.lines {
		current := l.text
		cont := false

		if i == 0 {
			current = multilineBeginRg.ReplaceAllString(current, `"`)
			cont = true
		} else {
			current = strings.Replace(current, s.multiLine.getIndent(), "", 1)
			current = multilineEndRg.ReplaceAllString(current, `"`)
		}

		if continueRg.MatchString(current) {
			cont = true
			current = continueRg.ReplaceAllString(current, "")
			current = regexp.MustCompile(`^(\s*)`).ReplaceAllString(current, "")
		}

		output += current

		if !cont && i != len(s.multiLine.lines)-1 && s.multiLine.lines[i+1].text != `"""` {
			output += "\n"
		}
	}

	s.multiLine = multi{lines: []line{}}

	return line{text: output, n: n}
}
