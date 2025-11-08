package variable

import "github.com/jorgefuertes/thenewquill/internal/adventure/db"

func (v Variable) Validate(allowNoID db.Allow) error {
	if err := v.ID.Validate(db.DontAllowSpecial); err != nil && !allowNoID {
		return err
	}

	if v.ID < db.MinMeaningfulID {
		return db.ErrInvalidLabelID
	}

	return nil
}

func (s Service) ValidateAll() error {
	for _, v := range s.All() {
		if err := v.Validate(db.DontAllowNoID); err != nil {
			return err
		}
	}

	return nil
}
