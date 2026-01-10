package database

import (
	"math"
	"sync"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
)

type (
	labels    map[uint32]string
	data      map[uint32]Record
	snapshots []data
)

type DB struct {
	mut         sync.Mutex
	lastLabelID uint32
	lastDataID  uint32
	labels      labels
	data        data
	frozen      bool
	snapshots   snapshots
	ram         data
}

func NewDB() *DB {
	db := new(DB)
	db.ResetDB()

	return db
}

// ResetDB resets the database to its initial empty state
func (db *DB) ResetDB() {
	db.mut.Lock()
	defer db.mut.Unlock()

	db.lastLabelID = 2
	db.lastDataID = 0
	db.labels = labels{0: "#undefined", 1: LabelAsterisk, 2: LabelUnderscore}
	db.data = make(data, 0)
	db.frozen = false
	db.snapshots = make(snapshots, 0)
}

func (db *DB) lock() {
	db.mut.Lock()
}

func (db *DB) unlock() {
	db.mut.Unlock()
}

func (db *DB) getNewDataID() (uint32, error) {
	if db.isFullOfRecords() {
		return 0, ErrDatabaseIsFull
	}

	db.lastDataID++
	return db.lastDataID, nil
}

func (db *DB) getNewLabelID() (uint32, error) {
	if db.isFullOfLabels() {
		return 0, ErrLabelsAreFull
	}

	db.lastLabelID++
	return db.lastLabelID, nil
}

func (db *DB) isFullOfRecords() bool {
	return len(db.data) == math.MaxUint32
}

func (db *DB) CountRecords() int {
	return len(db.data)
}

func (db *DB) CountRecordsByKind(k kind.Kind) int {
	db.lock()
	defer db.unlock()

	count := 0

	if k == kind.Label {
		return len(db.labels)
	}

	for _, r := range db.data {
		if r.Kind == k {
			count++
		}
	}

	return count
}
