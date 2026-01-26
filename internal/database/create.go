package database

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/jorgefuertes/thenewquill/internal/database/adapter"
)

func (db *DB) Create(entity adapter.Storeable) (uint32, error) {
	if db.IsFrozen() {
		return 0, ErrDatabaseIsFrozen
	}

	if db.isFullOfRecords() {
		return 0, ErrDatabaseIsFull
	}

	if entity.GetID() != 0 {
		return entity.GetID(), ErrIDFieldIsNotZero
	}

	if !db.ExistsLabelID(entity.GetLabelID()) {
		return entity.GetID(), ErrLabelNotFound
	}

	db.lock()
	defer db.unlock()

	r := Record{
		LabelID: entity.GetLabelID(),
		Kind:    entity.GetKind(),
		Data:    []byte{},
	}

	id, err := db.getNewDataID()
	if err != nil {
		return 0, err
	}

	entity.SetID(id)

	r.Data, err = cbor.Marshal(entity)
	if err != nil {
		return 0, err
	}

	db.data[id] = r

	return id, nil
}
