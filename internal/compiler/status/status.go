package status

import (
	"reflect"
	"slices"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
)

const stackSize = 5

type Status struct {
	Section          db.Kind
	Comment          line.Multi
	MultiLine        line.Multi
	Stack            []line.Line
	CurrentLabel     db.Label
	CurrentStoreable db.Storeable
	filenames        []string
}

func New() *Status {
	return &Status{
		Section:   db.None,
		Comment:   line.NewMulti(),
		MultiLine: line.NewMulti(),
		Stack:     []line.Line{},
	}
}

func (s *Status) PushFilename(filename string) {
	s.filenames = append(s.filenames, filename)
}

func (s *Status) PopFilename() {
	if len(s.filenames) == 0 {
		return
	}

	// remove the last element
	s.filenames = s.filenames[:len(s.filenames)-1]
}

func (s *Status) CurrentFilename() string {
	return s.filenames[len(s.filenames)-1]
}

// HasCurrentLabel returns true if there is a current label
func (s *Status) HasCurrentLabel() bool {
	return s.CurrentLabel.ID.IsDefined()
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

func (s *Status) AppendLine(l line.Line) {
	s.MultiLine.Append(l)
}

func (s *Status) Save(d *db.DB) error {
	if !s.HasCurrentLabel() || s.CurrentStoreable == nil {
		return nil
	}

	if _, err := d.Create(s.CurrentLabel.Name, s.CurrentStoreable); err != nil {
		return err
	}

	s.CurrentLabel = db.UndefinedLabel
	s.CurrentStoreable = nil

	return nil
}

func (s *Status) GetCurrentStoreable(dst any) bool {
	if s.CurrentStoreable == nil {
		return false
	}

	// dst must be a pointer
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		return false
	}

	dstValue.Elem().Set(reflect.ValueOf(s.CurrentStoreable))

	return true
}
