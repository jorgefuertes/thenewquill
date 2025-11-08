package status

import (
	"errors"
	"reflect"
	"slices"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
)

const stackSize = 5

type currentStoreable struct {
	label     db.Label
	storeable db.Storeable
	line      line.Line
	filename  string
}

type Status struct {
	db        *db.DB
	Section   kind.Kind
	Comment   line.Multi
	MultiLine line.Multi
	Stack     []line.Line
	current   *currentStoreable
	filenames []string
}

func New(d *db.DB) *Status {
	return &Status{
		db:        d,
		Section:   kind.None,
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
	if len(s.filenames) == 0 {
		return ""
	}

	return s.filenames[len(s.filenames)-1]
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

func (s *Status) SaveCurrentStoreable() cerr.CompilerError {
	if s.current == nil {
		return cerr.OK
	}

	_, err := s.db.Create(s.current.label.Name, s.current.storeable)
	if err != nil {
		return cerr.ErrDBCreate.WithStack(s.Stack).WithSection(s.Section).WithLine(s.current.line).
			WithFilename(s.current.filename).AddErr(err)
	}

	s.ClearCurrent()

	return cerr.OK
}

func (s *Status) GetCurrentStoreable(dst any) bool {
	if s.current == nil {
		return false
	}

	// dst must be a pointer
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		return false
	}

	dstValue.Elem().Set(reflect.ValueOf(s.current.storeable))

	return true
}

func (s *Status) SetCurrentStoreable(storeable db.Storeable) error {
	s.current.storeable = storeable

	return nil
}

func (s *Status) SetCurrentLabel(label db.Label) error {
	if s.current != nil {
		return errors.New("unexpected: cannot set a new label, current storeable already set")
	}

	s.current = &currentStoreable{
		label:    label,
		line:     s.Stack[len(s.Stack)-1],
		filename: s.CurrentFilename(),
	}

	return nil
}

func (s *Status) GetCurrentLabel() db.Label {
	if s.current == nil {
		return db.UndefinedLabel
	}

	return s.current.label
}

func (s *Status) ClearCurrent() {
	s.current = nil
}

func (s *Status) HasCurrent() bool {
	return s.current != nil
}
