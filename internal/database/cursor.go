package database

import "sync"

type cursor struct {
	mut  sync.Mutex
	i    int
	data data
}

func newCursor() *cursor {
	c := &cursor{}
	c.Reset()

	return c
}

func (c *cursor) Reset() {
	c.mut = sync.Mutex{}
	c.i = 0
	c.data = make(data, 0)
}

func (c *cursor) addOrReplace(id uint32, r Record) {
	c.data[id] = r
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
	checkDst(dst)

	_, r, ok := c.getByIndex(c.i)
	if !ok {
		return false
	}

	c.i++

	if err := r.Unmarshal(dst); err != nil {
		return false
	}

	return true
}

func (c *cursor) First(dst any) error {
	checkDst(dst)

	_, r, ok := c.getByIndex(0)
	if !ok {
		return ErrRecordNotFound
	}

	return r.Unmarshal(dst)
}

func (c *cursor) getFirstRecord() (uint32, Record, bool) {
	for id, r := range c.data {
		return id, r, true
	}

	return 0, Record{}, false
}

func (c *cursor) getByIndex(i int) (uint32, Record, bool) {
	if i >= len(c.data) {
		return 0, Record{}, false
	}

	var cur int
	for id, r := range c.data {
		if cur == i {
			return id, r, true
		}

		cur++
		if cur > i {
			break
		}
	}

	return 0, Record{}, false
}

func (db *DB) Query(filters ...filter) *cursor {
	c := newCursor()

	db.lock()
	defer db.unlock()

	for id, r := range db.data {
		if db.matchesAllFilters(r, filters...) {
			c.addOrReplace(id, r)
		}
	}

	if db.IsFrozen() {
		for _, snap := range db.snapshots {
			for id, r := range snap {
				if db.matchesAllFilters(r, filters...) {
					c.addOrReplace(id, r)
				}
			}
		}
	}

	return c
}
