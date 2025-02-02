package compiler

import (
	"slices"
)

const stackSize = 5

type status struct {
	section   section
	comment   multi
	multiLine multi
	stack     []line
}

func newStatus() *status {
	return &status{section: sectionNone}
}

func (s *status) appendStack(l line) {
	if len(s.stack) == stackSize {
		s.stack = slices.Delete(s.stack, 0, 1)
	}

	s.stack = append(s.stack, l)
}

func (s *status) setSection(section section, l line) {
	s.section = section
}

func (s *status) setComment(l line) {
	s.comment = multi{on: true, lines: []line{l}}
}

func (s *status) unsetComment() {
	s.comment = multi{}
}
