package db

import (
	"thenewquill/internal/compiler/section"
)

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

func (d *DB) Add(r Register) {
	d.regs = append(d.regs, r)
}

func (db *DB) GetHeaders() []string {
	return db.headers
}

func (db *DB) GetRegs() []Register {
	return db.regs
}

func (db *DB) GetRegsForSection(section section.Section) []Register {
	var regs []Register
	for _, r := range db.regs {
		if r.Section == section {
			regs = append(regs, r)
		}
	}

	return regs
}
