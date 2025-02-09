package compiler

import (
	"slices"

	"thenewquill/internal/adventure/loc"
)

const stackSize = 5

type seenLabel struct {
	label    string
	section  section
	line     line
	resolved bool
}

type status struct {
	section         section
	comment         multi
	multiLine       multi
	stack           []line
	labels          []seenLabel
	currentLocation loc.Location
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

func (s *status) appendLabel(label string, sec section, resolved bool, l line) {
	for _, la := range s.labels {
		if la.label == label && la.section == sec {
			return
		}
	}

	s.labels = append(s.labels, seenLabel{label: label, section: sec, resolved: resolved, line: l})
}

func (s *status) resolveLabel(label string, sec section) {
	for _, l := range s.labels {
		if l.label == label && l.section == sec {
			l.resolved = true
		}
	}
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

func (s *status) setCurrentLocation(label string) {
	s.currentLocation = loc.NewLocation(label, "", "")
}

func (s *status) unsetLocation() {
	s.currentLocation = loc.NewLocation("", "", "")
}

func (s *status) isCurrentLocation() bool {
	return s.currentLocation.Label != ""
}
