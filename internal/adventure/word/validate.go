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

		words2 := s.db.Query(db.FilterByKind(kind.Word), db.Filter("type", db.Equal, w.Type))
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
