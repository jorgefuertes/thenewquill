package database

import (
	"sync"
)

type cursor struct {
	mut  sync.Mutex
	i    int
	data []Record
}

func newCursor() *cursor {
	c := &cursor{}
	c.Reset()

	return c
}

func (c *cursor) Reset() {
	c.mut = sync.Mutex{}
	c.i = 0
	c.data = make([]Record, 0)
}

func (c *cursor) addOrReplace(r Record) {
	c.mut.Lock()
	defer c.mut.Unlock()

	c.data = append(c.data, r)
}

func (c *cursor) Exists() bool {
	return c.Count() > 0
}

func (c *cursor) Count() int {
	c.mut.Lock()
	defer c.mut.Unlock()

	return len(c.data)
}

func (c *cursor) Close() {
	c.Reset()
	c.data = nil
}

func (c *cursor) Next(dst any) bool {
	_, _ = checkEntity(dst)

	r, ok := c.getByIndex(c.i)
	if !ok {
		return false
	}

	if err := r.Unmarshal(dst); err != nil {
		return false
	}

	c.positionIncrement()

	return true
}

func (c *cursor) First(dst any) error {
	_, _ = checkEntity(dst)

	r, ok := c.getByIndex(0)
	if !ok {
		return ErrRecordNotFound
	}

	return r.Unmarshal(dst)
}

func (c *cursor) positionIncrement() {
	c.mut.Lock()
	defer c.mut.Unlock()

	c.i++
}

func (c *cursor) getByIndex(i int) (Record, bool) {
	c.mut.Lock()
	defer c.mut.Unlock()

	if i < 0 || i >= len(c.data) {
		return Record{}, false
	}

	return c.data[i], true
}

func (db *DB) Query(filters ...Filter) *cursor {
	c := newCursor()

	db.lock()
	defer db.unlock()

	for _, r := range db.data {
		if db.matchesAllFilters(r, filters...) {
			c.addOrReplace(r)
		}
	}

	if db.IsFrozen() {
		for _, snap := range db.snapshots {
			for _, r := range snap {
				if db.matchesAllFilters(r, filters...) {
					c.addOrReplace(r)
				}
			}
		}
	}

	return c
}
