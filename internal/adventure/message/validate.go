package message

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/internal/database/primitive"
)

func (m Message) Validate(allowNoID bool) error {
	if err := m.ID.ValidateID(false); err != nil && !allowNoID {
		return err
	}

	if m.ID < primitive.MinID && !allowNoID {
		return primitive.ErrInvalidID
	}

	if m.Text == "" {
		return ErrUndefinedText
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
		if err := m.Validate(false); err != nil {
			return err
		}
	}

	return nil
}
