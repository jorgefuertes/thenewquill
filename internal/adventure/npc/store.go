package npc

type Store []*NPC

func NewStore() Store {
	return Store{}
}

func (s Store) Get(label string) *NPC {
	for _, npc := range s {
		if npc.label == label {
			return npc
		}
	}

	return nil
}

func (s Store) Exists(label string) bool {
	for _, npc := range s {
		if npc.label == label {
			return true
		}
	}

	return false
}

// Set a new npc
func (s *Store) Set(n *NPC) error {
	if s.Exists(n.label) {
		return ErrNPCAlreadyExists
	}

	*s = append(*s, n)

	return nil
}
