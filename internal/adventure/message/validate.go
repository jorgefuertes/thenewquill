package message

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/pkg/validator"
)

func (m Message) Validate() error {
	if err := validator.Validate(m); err != nil {
		return err
	}

	if m.Plurals[One] != "" || m.Plurals[Many] != "" {
		if m.Plurals[One] == "" || m.Plurals[Many] == "" {
			return ErrUndefinedPlural
		}
	}

	return nil
}

func (s *Service) ValidateAll() error {
	msgs := s.db.Query(database.FilterByKind(kind.Message))
	defer msgs.Close()

	var m Message
	for msgs.Next(&m) {
		if err := m.Validate(); err != nil {
			return err
		}
	}

	return nil
}
