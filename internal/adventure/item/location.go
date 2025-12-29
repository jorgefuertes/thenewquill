package item

func (s *Service) IsAt(i Item, at uint32) bool {
	return i.At == at
}

func (s *Service) MoveTo(i *Item, to uint32) error {
	if s.IsContained(*i) {
		return ErrItemAlreadyContained
	}

	i.At = to

	return s.Update(i)
}
