package db

import (
	"reflect"

	"github.com/jorgefuertes/thenewquill/internal/adapter"
)

type cursor struct {
	i    int
	regs []adapter.Storeable
}

func (c *cursor) Reset() {
	c.i = 0
}

func (c *cursor) Count() int {
	return len(c.regs)
}

func (c *cursor) Close() {
	c.Reset()
	c.regs = nil
}

func (c *cursor) Next(dst any) bool {
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		return false
	}

	if c.i >= len(c.regs) {
		return false
	}

	c.i++

	dstValue.Elem().Set(reflect.ValueOf(c.regs[c.i-1]))

	return true
}

func (c *cursor) First(dst any) bool {
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		return false
	}

	if len(c.regs) == 0 {
		return false
	}

	dstValue.Elem().Set(reflect.ValueOf(c.regs[0]))

	return true
}

func (d *DB) Query(filters ...filter) *cursor {
	c := &cursor{i: 0, regs: make([]adapter.Storeable, 0)}

	for _, r := range d.Data {
		if matches(r, filters...) {
			c.regs = append(c.regs, r)
		}
	}

	return c
}
