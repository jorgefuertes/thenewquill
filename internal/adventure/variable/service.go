package variable

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Set(id uint32, value any) (uint32, error) {
	v, err := s.Get().WithID(id).First()
	if err != nil {
		return 0, err
	}

	v.SetValue(value)

	return v.ID, s.db.Update(v)
}

func (s *Service) SetByLabel(label string, value any) (uint32, error) {
	labelID, err := s.db.CreateLabel(label)
	if err != nil {
		return 0, err
	}

	v, err := s.Get().WithLabelID(labelID).First()
	if err == nil {
		v.SetValue(value)

		return v.ID, s.db.Update(v)
	}

	v.LabelID = labelID
	v.SetValue(value)

	return s.db.Create(v)
}

func (s *Service) Count() int {
	return s.db.CountRecordsByKind(kind.Variable)
}
