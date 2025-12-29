package config

import (
	"slices"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Set(label, v string) (uint32, error) {
	if !slices.Contains(allowedParamLabels, label) {
		return 0, ErrUnrecognizedConfigField
	}

	if s.Exists(label) {
		p, err := s.GetByLabel(label)
		if err != nil {
			return 0, err
		}

		p.V = v

		return p.ID, s.db.Update(p)
	}

	labelID, err := s.db.CreateLabel(label)
	if err != nil {
		return 0, err
	}

	p := &Param{LabelID: labelID, V: v}

	return s.db.Create(p)
}

func (s *Service) Get(id uint32) (*Param, error) {
	p := &Param{}
	err := s.db.Get(id, p)

	return p, err
}

func (s *Service) GetByLabel(label string) (*Param, error) {
	p := &Param{}
	err := s.db.GetByLabel(label, p)

	return p, err
}

func (s *Service) Exists(label string) bool {
	labelID, err := s.db.GetLabelID(label)
	if err != nil {
		return false
	}

	return s.db.Query(database.FilterByLabelID(labelID), database.FilterByKind(kind.Param)).Exists()
}

func (s *Service) Count() int {
	return s.db.CountRecordsByKind(kind.Param)
}

func (s *Service) GetParam(label string) string {
	c, err := s.GetByLabel(label)
	if err != nil {
		return ""
	}

	return c.V
}
