package db

import "reflect"

type cursor struct {
	i    int
	regs []Storeable
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

type filter struct {
	field string
	value any
}

func Filter(field string, value any) filter {
	return filter{field, value}
}

func (d *DB) Query(k Kind, filters ...filter) *cursor {
	c := &cursor{i: 0, regs: make([]Storeable, 0)}

	for _, r := range d.Data {
		if r.GetKind() != k {
			continue
		}

		match := true
		for _, f := range filters {
			field := reflect.ValueOf(r).FieldByName(f.field)

			if !field.IsValid() || field.Interface() != f.value {
				match = false
			}
		}

		if match {
			c.regs = append(c.regs, r)
		}
	}

	return c
}
