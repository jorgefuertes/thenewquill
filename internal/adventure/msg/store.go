package msg

type Store []*Msg

func NewStore() Store {
	return Store{}
}

func (s *Store) Set(m *Msg) error {
	if s.Exists(m.Type, m.Label) {
		return ErrMsgAlreadyExists
	}

	*s = append(*s, m)

	return nil
}

func (s Store) Exists(t MsgType, label string) bool {
	for _, msg := range s {
		if msg.Type == t && msg.Label == label {
			return true
		}
	}

	return false
}

func (s Store) Get(t MsgType, label string) *Msg {
	for _, m := range s {
		if m.Type == t && m.Label == label {
			return m
		}
	}

	return nil
}

func (s Store) Len() int {
	return len(s)
}
