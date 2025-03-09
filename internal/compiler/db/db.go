package db

const VERSION = "1"

type DB struct {
	headers []string
	regs    []Register
}

func NewDB() *DB {
	return &DB{
		headers: []string{"The New Quill Adventure Writing System"},
		regs:    make([]Register, 0),
	}
}

func (db *DB) From(i Exportable) {
	db.headers = append(db.headers, i.ExportHeaders()...)
	for sec, rows := range i.Export() {
		for _, row := range rows {
			db.regs = append(db.regs, Register{
				Section: sec,
				Fields:  row,
			})
		}
	}
}

func (db *DB) GetHeaders() []string {
	return db.headers
}

func (db *DB) GetRegs() []Register {
	return db.regs
}
