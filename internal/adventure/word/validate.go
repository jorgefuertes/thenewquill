package word

import (
	"errors"
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
	}

	// check for duplicates
	words2 := s.db.Query(db.FilterByKind(kind.Word))
	defer words2.Close()

	var w2 Word
	for words2.Next(&w2) {
		if w.GetID() == w2.GetID() {
			continue
		}

		for _, syn := range w.Synonyms {
			if w2.Is(w2.Type, syn) {
				return errors.Join(ErrDuplicatedWord, errors.New(s.db.GetLabelName(w.GetID())))
			}
		}
	}

	return nil
}
