package vars

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"

	"thenewquill/internal/compiler/rg"
)

type Store struct {
	lock *sync.Mutex
	Regs map[string]any
}

func NewStore() Store {
	return Store{
		lock: &sync.Mutex{},
		Regs: make(map[string]any, 0),
	}
}

func NewStoreFromMap(m map[string]any) Store {
	return Store{
		lock: &sync.Mutex{},
		Regs: m,
	}
}

func (s Store) Len() int {
	return len(s.Regs)
}

func (s *Store) Set(key string, value any) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Regs[key] = value
}

func (s *Store) SetFromString(key string, valueStr string) {
	if rg.Float.MatchString(valueStr) {
		value, _ := strconv.ParseFloat(valueStr, 64)
		s.Set(key, value)

		return
	}

	if rg.Int.MatchString(valueStr) {
		value, _ := strconv.ParseInt(valueStr, 10, 64)
		s.Set(key, int(value))

		return
	}

	if rg.Bool.MatchString(valueStr) {
		value, _ := strconv.ParseBool(valueStr)
		s.Set(key, value)

		return
	}

	s.Set(key, valueStr)
}

func (s *Store) SetAll(regs map[string]any) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Regs = regs
}

func (s *Store) Get(key string) any {
	s.lock.Lock()
	defer s.lock.Unlock()

	if v, ok := s.Regs[key]; ok {
		return v
	}

	return ""
}

func (s *Store) GetAll() map[string]any {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.Regs
}

func (s *Store) IsSet(key string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.Regs[key]

	return ok
}

func (s *Store) Unset(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.Regs, key)
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

func (s *Store) GetString(key string) string {
	v := s.Get(key)

	switch v := v.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%0.4f", v)
	case int:
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
