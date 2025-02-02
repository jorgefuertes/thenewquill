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
	on    bool
	lines []line
}

func (m *multi) append(l line) {
	m.lines = append(m.lines, l)
}

func (m multi) isOn() bool {
	return m.on
}

func (s *status) startMultiLine(l line) {
	s.multiLine = multi{on: true, lines: []line{l}}
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
	l := s.multiLine.lines[0]
	l.text = multilineBeginRg.ReplaceAllString(l.text, `"`)
	for i, l2 := range s.multiLine.lines[1:] {
		l2.text = strings.Replace(l2.text, s.multiLine.getIndent(), "", 1)
		if i == 0 {
			l.text += l2.text
			continue
		}

		if continueRg.MatchString(l.text) {
			l.text = continueRg.ReplaceAllString(l.text, "")
		} else {
			l.text += "\n"
		}

		l.text += l2.text
	}

	l.text += `"`
	s.multiLine = multi{on: false, lines: []line{}}

	return l
}
