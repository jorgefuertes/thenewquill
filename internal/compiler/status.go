package compiler

type status struct {
	section   section
	inComment bool
	stack     []line
	inMulti   bool
	multiLine line
}

func newStatus() *status {
	return &status{section: sectionNone}
}

func (s *status) setSection(section section, l line) {
	s.section = section
	s.stack = append(s.stack, l)
}

func (s *status) setComment(l line) {
	s.inComment = true
	s.stack = append(s.stack, l)
}

func (s *status) unsetComment() {
	s.inComment = false
	if len(s.stack) > 0 {
		s.stack = s.stack[:len(s.stack)-1]
	}
}

func (s *status) getLastLine() line {
	if len(s.stack) > 0 {
		return s.stack[len(s.stack)-1]
	}

	return line{}
}

func (s *status) setMultiLine(l line) {
	s.inMulti = true
	s.multiLine = l
}

func (s *status) unsetMultiLine() {
	s.inMulti = false
	s.multiLine = line{}
}

func (s *status) getMultiLine() line {
	return s.multiLine
}

func (s *status) appendMultiLine(l line) {
	s.multiLine.text += l.text
}
