package character

func (s *Store) Validate() error {
	if s.GetHuman() == nil {
		return ErrNoHuman
	}

	s.mut.Lock()
	defer s.mut.Unlock()

	humans := 0
	for _, p := range s.chars {
		if p.Human {
			humans++
		}
	}

	if humans > 1 {
		return ErrOnlyOneHuman
	}

	return nil
}
