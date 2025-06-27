package message

import "github.com/jorgefuertes/thenewquill/internal/adventure/db"

func (m Message) Validate(allowNoID db.Allow) error {
	if err := m.ID.Validate(db.DontAllowSpecial); err != nil && !allowNoID {
		return err
	}

	if m.ID < db.MinMeaningfulID && !allowNoID {
		return db.ErrInvalidLabelID
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
	msgs := s.db.Query(db.Messages)
	defer msgs.Close()

	var m Message
	for msgs.Next(&m) {
		if err := m.Validate(db.DontAllowNoID); err != nil {
			return err
		}
	}

	return nil
}
