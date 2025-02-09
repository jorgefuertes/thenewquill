package store

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

func (d *Store) Set(key string, value any) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.regs[key] = value
}

func (d *Store) Get(key string) any {
	d.lock.Lock()
	defer d.lock.Unlock()

	if v, ok := d.regs[key]; ok {
		return v
	}

	return ""
}

func (d *Store) IsSet(key string) bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	_, ok := d.regs[key]

	return ok
}

func (d *Store) Unset(key string) {
	d.lock.Lock()
	defer d.lock.Unlock()

	delete(d.regs, key)
}

func (d *Store) GetBool(key string) bool {
	v := d.Get(key)

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

func (d *Store) GetInt(key string) int {
	v := d.Get(key)

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

func (d *Store) GetFloat(key string) float64 {
	v := d.Get(key)

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

func (d *Store) Count() int {
	return len(d.regs)
}
