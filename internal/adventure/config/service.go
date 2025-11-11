package config

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

type Service struct {
	db *db.DB
}

func NewService(d *db.DB) *Service {
	return &Service{db: d}
}

func (s *Service) Set(field string, v string) error {
	if !isKeyAllowed(field) {
		return ErrUnrecognizedConfigField
	}

	l, err := s.db.AddLabel(field)
	if err != nil {
		return err
	}

	c := Param{ID: l.ID, V: v}
	if err := s.db.Append(c); err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(v Param) error {
	return s.db.Update(v)
}

func (s *Service) Get(id id.ID) (Param, error) {
	v := Param{}
	err := s.db.Get(id, &v)

	return v, err
}

func (s *Service) GetField(name string) string {
	v := Param{}
	if err := s.db.GetByLabel(name, &v); err != nil {
		return ""
	}

	return v.V
}

func (s *Service) All() []Param {
	values := make([]Param, 0)

	q := s.db.Query(db.FilterByKind(kind.Param))
	var value Param
	for q.Next(&value) {
		values = append(values, value)
	}

	return values
}

func (s *Service) Count() int {
	return s.db.CountByKind(kind.Param)
}
