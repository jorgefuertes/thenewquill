package database

func (db *DB) Freeze() {
	db.lock()
	defer db.unlock()

	db.frozen = true
}

func (db *DB) IsFrozen() bool {
	return db.frozen
}

func (db *DB) Snapshot() {
	if !db.IsFrozen() {
		return
	}

	db.lock()
	defer db.unlock()

	db.snapshots = append(db.snapshots, make(data, 0))
}

func (db *DB) SnapBack() bool {
	if !db.IsFrozen() {
		return false
	}

	db.lock()
	defer db.unlock()

	if len(db.snapshots) == 0 {
		return false
	}

	db.snapshots = db.snapshots[:len(db.snapshots)-1]

	return true
}
