package config

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
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

	if !s.db.ExistsLabelName(field) {
		if _, err := s.db.AddLabel(field, false); err != nil {
			return err
		}
	}

	label, err := s.db.GetLabelByName(field)
	if err != nil {
		return err
	}

	c := Value{ID: label.ID, V: v}
	if err := s.db.Append(c); err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(v Value) error {
	return s.db.Update(v)
}

func (s *Service) Get(id db.ID) (Value, error) {
	v := Value{}
	err := s.db.Get(id, &v)

	return v, err
}

func (s *Service) GetField(name string) string {
	v := Value{}
	if err := s.db.GetByLabel(name, &v); err != nil {
		return ""
	}

	return v.V
}

func (s *Service) All() []Value {
	values := make([]Value, 0)

	q := s.db.Query(db.FilterByKind(kind.Config))
	var value Value
	for q.Next(&value) {
		values = append(values, value)
	}

	return values
}

func (s *Service) Count() int {
	return s.db.CountByKind(kind.Config)
}
