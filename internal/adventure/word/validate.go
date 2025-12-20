package word

import (
	"errors"
	"fmt"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/internal/database/primitive"
)

func (w Word) Validate(allowNoID bool) error {
	if w.ID == primitive.UndefinedID && !allowNoID {
		return nil
	}

	if err := w.ID.ValidateID(false); err != nil && !allowNoID {
		return fmt.Errorf("ID %q: %w", w.ID, err)
	}

	if len(w.Synonyms) == 0 {
		return ErrEmptyWord
	}

	return nil
}

func (s *Service) ValidateAll() error {
	words := s.db.Query(database.FilterByKind(kind.Word))
	defer words.Close()

	var w Word
	for words.Next(&w) {
		if err := w.Validate(false); err != nil {
			return err
		}

		words2 := s.db.Query(database.FilterByKind(kind.Word), database.Filter("type", database.Equal, w.Type))
		defer words2.Close()

		var w2 Word
		for words2.Next(&w2) {
			if w.ID == w2.ID {
				continue
			}

			for _, syn := range w.Synonyms {
				if w2.HasSynonym(syn) {
					return errors.Join(
						ErrDuplicatedWord,
						fmt.Errorf(
							"duplicated synonym %q between %s %q and %s %q",
							syn,
							w.Type,
							w.Synonyms[0],
							w2.Type,
							w2.Synonyms[0],
						),
					)
				}
			}
		}
	}

	return nil
}
