package character

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

func ValidateAll(d *db.DB) error {
	humans := 0

	chars := make([]Character, 0)
	if err := d.GetByKindAs(db.Chars, db.NoSubKind, &chars); err != nil {
		return err
	}

	for _, c := range chars {
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
