package database

import (
	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
)

type testStoreable struct {
	ID      primitive.ID
	LabelID primitive.ID
	Title   string
	At      primitive.ID
	OK      bool
	NOOK    bool
	Weight  int16
	Names   []string
	Numbers []int
}

var _ adapter.Storeable = &testStoreable{}

func (s testStoreable) GetID() primitive.ID {
	return s.ID
}

func (s *testStoreable) SetID(id primitive.ID) {
	s.ID = id
}

func (s testStoreable) GetLabelID() primitive.ID {
	return s.LabelID
}

func (s *testStoreable) SetLabelID(id primitive.ID) {
	s.LabelID = id
}

func (s testStoreable) Validate(allowNoID bool) error {
	return nil
}
