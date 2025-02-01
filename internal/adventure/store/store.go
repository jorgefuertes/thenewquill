package store

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"sync"
)

type Store struct {
	lock *sync.Mutex
	regs map[string]string
}

func New() Store {
	return Store{
		lock: &sync.Mutex{},
		regs: make(map[string]string, 0),
	}
}

func (d *Store) Set(key string, value any) {
	d.lock.Lock()
	defer d.lock.Unlock()

	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		d.regs[key] = value.(string)
	case reflect.Int:
		d.regs[key] = fmt.Sprintf("%d", value.(int))
	case reflect.Float64:
		d.regs[key] = fmt.Sprintf("%.2f", value.(float64))
	case reflect.Bool:
		d.regs[key] = fmt.Sprintf("%t", value.(bool))
	default:
		d.regs[key] = fmt.Sprintf("%v", value)
	}
}

func (d *Store) Get(key string) string {
	d.lock.Lock()
	defer d.lock.Unlock()

	if v, ok := d.regs[key]; ok {
		return v
	}

	return ""
}

func (d *Store) IsSet(key string) bool {
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
	re := regexp.MustCompile(`(?i)^\s*(true|yes|1|sí|si|ok|on)\s*$`)

	return re.MatchString(v)
}

func (d *Store) GetNumber(key string) float64 {
	v := d.Get(key)
	if v == "" {
		return 0
	}

	if v == "true" {
		return 1
	}

	if v == "false" {
		return 0
	}

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0
	}

	return f
}
