package character

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

func (c Character) Validate(allowNoID db.Allow) error {
	if err := c.ID.Validate(db.DontAllowSpecial); err != nil && !allowNoID {
		return err
	}

	if c.Description == "" {
		return ErrEmptyDescription
	}

	if err := c.NounID.Validate(db.DontAllowSpecial); err != nil {
		return fmt.Errorf("name ID %q: %w", c.NounID, err)
	}

	if err := c.AdjectiveID.Validate(db.AllowSpecial); err != nil {
		return fmt.Errorf("adjective ID %q: %w", c.AdjectiveID, err)
	}

	if err := c.LocationID.Validate(db.DontAllowSpecial); err != nil {
		return fmt.Errorf("location ID %q: %w", c.LocationID, err)
	}

	return nil
}

func (s *Service) ValidateAll() error {
	humans := 0

	chars := s.db.Query(db.FilterByKind(db.Characters))
	defer chars.Close()

	c := Character{}
	for chars.Next(&c) {
		if err := c.Validate(db.DontAllowNoID); err != nil {
			return err
		}

		if c.Human {
			humans++
		}
	}

	if humans > 1 {
		return ErrOnlyOneHuman
	}

	if humans == 0 {
		return ErrNoHuman
	}

	return nil
}
