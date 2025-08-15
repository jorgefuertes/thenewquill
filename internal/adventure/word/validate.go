package word

import (
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

func (w Word) Validate(allowNoID db.Allow) error {
	if w.ID == db.UndefinedLabel.ID && allowNoID == db.AllowNoID {
		return nil
	}

	if err := w.ID.Validate(db.DontAllowSpecial); err != nil && !allowNoID {
		return fmt.Errorf("ID %q: %w", w.ID, err)
	}

	if len(w.Synonyms) == 0 {
		return ErrEmptyWord
	}

	return nil
}

func (s *Service) ValidateAll() error {
	words := s.db.Query(db.FilterByKind(kind.Word))
	defer words.Close()

	var w Word
	for words.Next(&w) {
		if err := w.Validate(db.DontAllowNoID); err != nil {
			return err
		}

		// check for duplicates
	}

	return nil
}
