package character

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/pkg/validator"
)

func (s *Service) ValidateAll() error {
	humansCount := s.db.Query(database.FilterByKind(kind.Character), database.NewFilter("Human", database.Equal, true)).
		Count()
	if humansCount == 0 {
		return ErrNoHuman
	} else if humansCount > 1 {
		return ErrOnlyOneHuman
	}

	chars := s.db.Query(database.FilterByKind(kind.Character))
	defer chars.Close()

	c := &Character{}
	for chars.Next(c) {
		if err := validator.Validate(c); err != nil {
			return err
		}
	}

	return nil
}
