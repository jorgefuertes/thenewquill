package database

import "github.com/jorgefuertes/thenewquill/internal/adventure/kind"

func (db *DB) Exists(id uint32) bool {
	db.lock()
	defer db.unlock()

	_, ok := db.data[id]

	return ok
}

func (db *DB) Get(id uint32, dest any) error {
	k := kind.KindOf(dest)
	c := db.Query(FilterByID(id), FilterByKind(k))

	return c.First(dest)
}

func (db *DB) GetByLabel(label string, dest any) error {
	k := kind.KindOf(dest)
	c := db.Query(FilterByLabel(label), FilterByKind(k))

	return c.First(dest)
}

func (db *DB) GetByLabelID(id uint32, dest any) error {
	k := kind.KindOf(dest)
	c := db.Query(FilterByLabelID(id), FilterByKind(k))

	return c.First(dest)
}
