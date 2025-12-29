package database

import (
	"encoding/json"

	"github.com/fxamacker/cbor/v2"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/pkg/log"
)

func (db *DB) Create(entity any) (uint32, error) {
	if db.IsFrozen() {
		return 0, ErrDatabaseIsFrozen
	}

	if db.isFullOfRecords() {
		return 0, ErrDatabaseIsFull
	}

	id, labelID := checkEntity(entity)

	j, _ := json.Marshal(entity)
	log.Debug("🗄️ [DB] Create: %T:%s", entity, j)

	if id != 0 {
		return id, ErrIDFieldIsNotZero
	}

	if !db.ExistsLabelID(labelID) {
		return id, ErrLabelNotFound
	}

	db.lock()
	defer db.unlock()

	r := Record{
		LabelID: labelID,
		Kind:    kind.KindOf(entity),
		Data:    []byte{},
	}

	id, err := db.getNewDataID()
	if err != nil {
		return 0, err
	}

	setID(entity, id)

	r.Data, err = cbor.Marshal(entity)
	if err != nil {
		return 0, err
	}

	db.data[id] = r

	return id, nil
}
