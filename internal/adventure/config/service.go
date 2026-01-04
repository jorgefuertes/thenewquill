package config

import (
	"slices"

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

	if s.Get().WithLabel(label).Exists() {
		p, err := s.Get().WithLabel(label).First()
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

func (s *Service) GetValueOrBlank(label string) string {
	p, err := s.Get().WithLabel(label).First()
	if err != nil {
		return ""
	}

	return p.V
}
