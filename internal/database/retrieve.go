package database

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database/adapter"
)

func (db *DB) Get(id uint32, dest adapter.Storeable) error {
	k := kind.KindOf(dest)
	c := db.Query(FilterByID(id), FilterByKind(k))

	return c.First(dest)
}

func (db *DB) GetByLabel(label string, dest adapter.Storeable) error {
	k := kind.KindOf(dest)
	c := db.Query(FilterByLabel(label), FilterByKind(k))

	return c.First(dest)
}

func (db *DB) GetByLabelID(id uint32, dest adapter.Storeable) error {
	k := kind.KindOf(dest)
	c := db.Query(FilterByLabelID(id), FilterByKind(k))

	return c.First(dest)
}

func (db *DB) GetKind(id uint32) kind.Kind {
	db.lock()
	defer db.unlock()

	if r, ok := db.data[id]; ok {
		return r.Kind
	}

	return kind.None
}

func (db *DB) GetKindByLabelID(labelID uint32) kind.Kind {
	db.lock()
	defer db.unlock()

	for _, r := range db.data {
		if r.LabelID == labelID {
			return r.Kind
		}
	}

	return kind.None
}
