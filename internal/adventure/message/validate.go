package message

import "github.com/jorgefuertes/thenewquill/internal/adventure/db"

func (m Message) Validate(allowNoID db.Allow) error {
	if err := m.ID.Validate(db.DontAllowSpecial); err != nil && !allowNoID {
		return err
	}

	if m.ID < db.MinMeaningfulID {
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
	var m Message
	for s.db.Query(db.Messages).Next(&m) {
		if err := m.Validate(db.DontAllowNoID); err != nil {
			return err
		}
	}

	return nil
}
