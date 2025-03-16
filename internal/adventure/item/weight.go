package item

func (i *Item) TotalWeight(s *Store) int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.recursiveWeight(i)
}

func (s *Store) recursiveWeight(i *Item) int {
	if !i.IsContainer {
		return i.Weight
	}

	w := i.Weight
	for _, f := range s.items {
		if f.Inside == i {
			w += s.recursiveWeight(f)
		}
	}

	return w
}
