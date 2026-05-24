package database

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
