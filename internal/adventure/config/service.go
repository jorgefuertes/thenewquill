package config

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/internal/database/primitive"
)

type Service struct {
	db *database.DB
}

func NewService(db *database.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Set(labelOrString any, v string) (uint32, error) {
	label, err := primitive.LabelFromLabelOrString(labelOrString)
	if err != nil {
		return primitive.UndefinedID, err
	}

	if !isLabelAllowed(label) {
		return primitive.UndefinedID, ErrUnrecognizedConfigField
	}

	labelID, err := s.db.CreateLabelIfNotExists(label, database.DenyCompositeLabel)
	if err != nil {
		return primitive.UndefinedID, err
	}

	if s.db.Exists(database.FilterByKind(kind.Param), database.FilterByLabelID(labelID)) {
		p, err := s.GetByLabel(label)
		if err != nil {
			return primitive.UndefinedID, err
		}

		p.V = v

		return p.ID, s.db.Update(p)
	}

	p := New(primitive.UndefinedID, labelID, v)

	return s.db.Create(p)
}

func (s *Service) Get(id uint32) (*Param, error) {
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
