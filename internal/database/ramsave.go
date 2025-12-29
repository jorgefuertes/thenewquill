package database

func (db *DB) RamSave() {
	db.lock()
	defer db.unlock()

	db.ram = make(data, 0)

	for _, snap := range db.snapshots {
		for id, r := range snap {
			db.ram[id] = r
		}
	}
}

func (db *DB) RamLoad() {
	db.lock()
	defer db.unlock()

	db.snapshots = append(db.snapshots, db.ram)
}
