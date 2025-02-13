package compiler

import (
	"slices"
)

const stackSize = 5

type undef struct {
	section section
	line    line
	file    string
	label   string
}

type status struct {
	section      section
	comment      multi
	multiLine    multi
	stack        []line
	currentLabel string
	undef        []undef
}

func newStatus() *status {
	return &status{section: sectionNone}
}

func (s *status) hasCurrentLabel() bool {
	return s.currentLabel != ""
}

func (s *status) setCurrentLabel(label string) {
	s.currentLabel = label
}

func (s *status) unsetCurrentLabel() {
	s.currentLabel = ""
}

func (s *status) appendStack(l line) {
	if len(s.stack) == stackSize {
		s.stack = slices.Delete(s.stack, 0, 1)
	}

	s.stack = append(s.stack, l)
}

func (s *status) setSection(section section) {
	s.section = section
}

func (s *status) setComment(l line) {
	s.comment = multi{lines: []line{l}}
}

func (s *status) unsetComment() {
	s.comment = multi{}
}

func (s *status) setUndef(label string, section section, l line, file string) {
	s.undef = append(s.undef, undef{section: section, line: l, file: file, label: label})
}

func (s *status) setDef(label string, section section) {
	for i, u := range s.undef {
		if u.label == label && u.section == section {
			s.undef = append(s.undef[:i], s.undef[i+1:]...)
		}
	}
}
