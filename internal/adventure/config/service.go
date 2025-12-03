package config

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/database"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Set(labelOrString any, v string) (primitive.ID, error) {
	label, err := primitive.LabelFromLabelOrString(labelOrString)
	if err != nil {
		return primitive.UndefinedID, err
	}

	if !isLabelAllowed(label) {
		return primitive.UndefinedID, ErrUnrecognizedConfigField
	}

	if s.Exists(label) {
		p, err := s.GetByLabel(label)
		if err != nil {
			return primitive.UndefinedID, err
		}

		p.V = v

		return p.ID, s.db.Update(p)
	}

	labelID, err := s.db.CreateLabelIfNotExists(label, false)
	if err != nil {
		return primitive.UndefinedID, err
	}

	p := New(primitive.UndefinedID, labelID, v)

	return s.db.Create(p)
}

func (s *Service) Get(id primitive.ID) (*Param, error) {
	p := &Param{}
	err := s.db.Get(id, p)

	return p, err
}

func (s *Service) GetByLabel(labelOrString any) (*Param, error) {
	label, err := primitive.LabelFromLabelOrString(labelOrString)
	if err != nil {
		return nil, err
	}

	p := &Param{}
	err = s.db.GetByLabel(label, p)

	return p, err
}

func (s *Service) Exists(labelOrString any) bool {
	label, err := primitive.LabelFromLabelOrString(labelOrString)
	if err != nil {
		return false
	}

	labelID, err := s.db.GetLabelID(label)
	if err != nil {
		return false
	}

	return s.db.Exists(database.FilterByLabelID(labelID), database.FilterByKind(kind.Param))
}

func (s *Service) Count() int {
	return s.db.Count(database.FilterByKind(kind.Param))
}

func (s *Service) GetField(labelOrString any) string {
	label, err := primitive.LabelFromLabelOrString(labelOrString)
	if err != nil {
		return ""
	}

	c, err := s.GetByLabel(label)
	if err != nil {
		return ""
	}

	return c.V
}
