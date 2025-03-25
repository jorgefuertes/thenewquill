package item

import "errors"

func (s *Store) Validate() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, item := range s.items {
		if err := item.validate(); err != nil {
			return errors.Join(ErrItemValidationFailed, err, errors.New(item.Label))
		}
	}

	for i, item := range s.items {
		for i2, item2 := range s.items {
			if i == i2 {
				continue
			}

			if i != i2 && item.Label == item2.Label {
				return errors.Join(ErrDuplicatedItemLabel, errors.New(item.Label))
			}

			if item.Noun.Label == item2.Noun.Label && item.Adjective.Label == item2.Adjective.Label {
				return errors.Join(ErrDuplicatedNounAdj, errors.New(item.Label))
			}
		}
	}

	return nil
}
