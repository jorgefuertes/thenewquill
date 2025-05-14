package words

import (
	"errors"
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

func (w Word) Validate() error {
	if w.ID < db.UndefinedLabel.ID {
		return ErrEmptyLabel
	}

	if w.Type == db.NoSubKind && w.ID != db.UndefinedLabel.ID {
		return errors.Join(ErrUnknownWordType, fmt.Errorf("ID: '%d'", w.ID))
	}

	return nil
}

func ValidateAll(d *db.DB) error {
	words := make([]Word, 0)

	if err := d.GetByKindAs(db.Words, db.NoSubKind, &words); err != nil {
		return err
	}

	for _, w := range words {
		if err := w.Validate(); err != nil {
			return err
		}

		// check for duplicates
		for _, w2 := range words {
			if w.GetID() == w2.GetID() {
				continue
			}

			for _, s := range w.Synonyms {
				if w2.IsSynonymAndType(s, w2.Type) {
					return errors.Join(ErrDuplicatedWord, errors.New(d.GetLabelName(w.GetID())))
				}
			}
		}
	}

	return nil
}
