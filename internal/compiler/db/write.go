package db

import (
	"compress/zlib"
	"fmt"
	"io"

	"thenewquill/internal/compiler/section"
	"thenewquill/internal/util"
)

const (
	fieldSep     byte = 0x1F
	regSep       byte = 0x1E
	beginRecords byte = 0x1D
	endRecords   byte = 0x1C
)

func (db *DB) Reset() {
	db.headers = make([]string, 0)
	db.regs = make([]Register, 0)
}

func (db *DB) AddHeader(lines ...string) {
	db.headers = append(db.headers, lines...)
}

func (db *DB) AddReg(s section.Section, fields ...any) {
	r := Register{Section: s, Fields: []string{}}

	for _, f := range fields {
		r.Fields = append(r.Fields, util.ValueToString(f))
	}

	db.regs = append(db.regs, r)
}

func (db *DB) Write(w io.Writer) error {
	// write plain headers
	for _, h := range db.headers {
		if _, err := w.Write([]byte(h + "\n")); err != nil {
			return err
		}
	}

	// write begin records
	if _, err := w.Write([]byte{beginRecords}); err != nil {
		return err
	}

	// write zlib-compressed registers
	zw := zlib.NewWriter(w)

	for _, r := range db.regs {
		err := writeReg(zw, r)
		if err != nil {
			return err
		}
	}

	if err := zw.Flush(); err != nil {
		return err
	}
	zw.Close()

	// write end records
	if _, err := w.Write([]byte{endRecords}); err != nil {
		return err
	}

	// hash
	if _, err := w.Write(db.Hash()); err != nil {
		return err
	}

	return nil
}

func writeReg(zw *zlib.Writer, r Register) error {
	_, err := zw.Write([]byte(fmt.Sprintf("%d", r.Section)))
	if err != nil {
		return err
	}

	_, err = zw.Write([]byte{fieldSep})
	if err != nil {
		return err
	}

	for _, f := range r.Fields {
		_, err := zw.Write([]byte(f))
		if err != nil {
			return err
		}

		_, err = zw.Write([]byte{fieldSep})
		if err != nil {
			return err
		}
	}

	if _, err := zw.Write([]byte{regSep}); err != nil {
		return err
	}

	return nil
}
