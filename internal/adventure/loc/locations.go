package loc

type Locations struct {
	locations []*Location
}

func New() Locations {
	return Locations{locations: make([]*Location, 0)}
}

func (s *Locations) Get(label string) *Location {
	for _, l := range s.locations {
		if l.Label == label {
			return l
		}
	}

	return nil
}

func (s *Locations) IndexOf(label string) int {
	for i, l := range s.locations {
		if l.Label == label {
			return i
		}
	}

	return -1
}

// Set a new location
// overwrites any existing location with the same label
func (s *Locations) Set(label, title, desc string) *Location {
	l := &Location{
		Label:       label,
		Title:       title,
		Description: desc,
		Conns:       make([]Connection, 0),
	}

	oldIdx := s.IndexOf(l.Label)
	if oldIdx != -1 {
		s.locations[oldIdx] = l

		return l
	}

	s.locations = append(s.locations, l)

	return l
}

func (s *Locations) Exists(label string) bool {
	for _, l := range s.locations {
		if l.Label == label {
			return true
		}
	}

	return false
}
