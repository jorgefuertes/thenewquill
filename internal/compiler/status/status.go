package status

import (
	"slices"

	"thenewquill/internal/compiler/line"
	"thenewquill/internal/compiler/section"
)

const stackSize = 5

type Undef struct {
	Section section.Section
	Line    line.Line
	File    string
	Label   string
}

type Status struct {
	Section      section.Section
	Comment      line.Multi
	MultiLine    line.Multi
	Stack        []line.Line
	CurrentLabel string
	Undefs       []Undef
	filenames    []string
}

func New() *Status {
	return &Status{
		Section:   section.None,
		Comment:   line.NewMulti(),
		MultiLine: line.NewMulti(),
		Stack:     []line.Line{},
		Undefs:    []Undef{},
	}
}

func (s *Status) PushFilename(filename string) {
	s.filenames = append(s.filenames, filename)
}

func (s *Status) PopFilename() {
	if len(s.filenames) == 0 {
		return
	}

	s.filenames = slices.Delete(s.filenames, 0, 1)
}

func (s *Status) CurrentFilename() string {
	return s.filenames[len(s.filenames)-1]
}

// HasCurrentLabel returns true if there is a current label
func (s *Status) HasCurrentLabel() bool {
	return s.CurrentLabel != ""
}

func (s *Status) AppendStack(l line.Line) {
	if len(s.Stack) == stackSize {
		s.Stack = slices.Delete(s.Stack, 0, 1)
	}

	s.Stack = append(s.Stack, l)
}

func (s *Status) SetComment(l line.Line) {
	s.Comment = line.NewMulti(l)
}

func (s *Status) UnsetComment() {
	s.Comment.Clear()
}

func (s *Status) SetUndef(label string, section section.Section, l line.Line) {
	s.Undefs = append(s.Undefs, Undef{Section: section, Line: l, File: s.CurrentFilename(), Label: label})
}

func (s *Status) SetDef(label string, section section.Section) {
	for i, u := range s.Undefs {
		if u.Label == label && u.Section == section {
			s.Undefs = append(s.Undefs[:i], s.Undefs[i+1:]...)
		}
	}
}

func (s *Status) IsUndef(label string, section section.Section) bool {
	for _, u := range s.Undefs {
		if u.Label == label && u.Section == section {
			return true
		}
	}

	return false
}

func (s *Status) HasAnyUndef() bool {
	return len(s.Undefs) > 0
}

func (s *Status) AppendLine(l line.Line) {
	s.MultiLine.Append(l)
}
