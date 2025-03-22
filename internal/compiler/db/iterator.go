package db

import (
	"sync"

	"thenewquill/internal/compiler/section"
)

type iterator struct {
	mut   *sync.Mutex
	db    *DB
	sec   section.Section
	index int
}

func (d *DB) NewIterator(sec section.Section) *iterator {
	return &iterator{mut: &sync.Mutex{}, db: d, sec: sec, index: 0}
}

func (i *iterator) Next() *Record {
	i.mut.Lock()
	defer i.mut.Unlock()

	if i.index >= len(i.db.Records) {
		return nil
	}

	for i.index < len(i.db.Records) {
		i.index++

		if i.db.Records[i.index].Section == i.sec {
			r := i.db.Records[i.index]

			return &r
		}
	}

	return nil
}

func (i *iterator) Reset() {
	i.mut.Lock()
	defer i.mut.Unlock()

	i.index = 0
}
