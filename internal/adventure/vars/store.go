package vars

import (
	"regexp"
	"strconv"
	"sync"
)

type Store struct {
	lock *sync.Mutex
	regs map[string]any
}

func New() Store {
	return Store{
		lock: &sync.Mutex{},
		regs: make(map[string]any, 0),
	}
}

func (s *Store) All() map[string]any {
	return s.regs
}

func (s *Store) Set(key string, value any) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.regs[key] = value
}

func (s *Store) Get(key string) any {
	s.lock.Lock()
	defer s.lock.Unlock()

	if v, ok := s.regs[key]; ok {
		return v
	}

	return ""
}

func (s *Store) IsSet(key string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.regs[key]

	return ok
}

func (s *Store) Unset(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.regs, key)
}

func (s *Store) GetBool(key string) bool {
	v := s.Get(key)

	switch v := v.(type) {
	case string:
		re := regexp.MustCompile(`(?i)^\s*(true|yes|1|sí|si|ok|on)\s*$`)

		return re.MatchString(v)
	case bool:
		return v
	case int:
		return v > 0
	case float64:
		return v > 0
	default:
		return false
	}
}

func (s *Store) GetInt(key string) int {
	v := s.Get(key)

	switch v := v.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return int(i)
		}

		return 0
	case bool:
		if v {
			return 1
		}

		return 0
	default:
		return 0
	}
}

func (s *Store) GetFloat(key string) float64 {
	v := s.Get(key)

	switch v := v.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return float64(i)
		}

		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0
		}

		return f
	case bool:
		if v {
			return 1.0
		}

		return 0.0
	default:
		return 0.0
	}
}

func (s *Store) Count() int {
	return len(s.regs)
}
