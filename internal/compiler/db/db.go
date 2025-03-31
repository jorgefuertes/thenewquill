package db

import (
	"crypto/sha512"
	"encoding/gob"
	"fmt"
	"strings"

	"thenewquill/internal/compiler/section"
)

const (
	EOH byte = 0x1C // End of Headers
	EOR byte = 0x1D // End of Records
)

type DB struct {
	Headers []string
	Records []Record
}

func New() *DB {
	gob.Register(map[string]string{})
	gob.Register(map[string]any{})

	return &DB{
		Headers: make([]string, 0),
		Records: make([]Record, 0),
	}
}

func (db *DB) Reset() {
	db.Headers = make([]string, 0)
	db.Records = make([]Record, 0)
}

func (db *DB) Hash() ([]byte, error) {
	hasher := sha512.New()
	_, err := hasher.Write([]byte(strings.Join(db.Headers, "")))
	if err != nil {
		return nil, err
	}

	for _, r := range db.Records {
		// section
		_, err := hasher.Write([]byte{r.Section.Byte()})
		if err != nil {
			return nil, err
		}

		// label
		_, err = hasher.Write([]byte(r.Label))
		if err != nil {
			return nil, err
		}

		// fields
		for _, f := range r.Fields {
			_, err := fmt.Fprintf(hasher, "%v", f)
			if err != nil {
				return nil, err
			}
		}
	}

	return hasher.Sum(nil), nil
}

func (db *DB) AddHeader(lines ...string) {
	db.Headers = append(db.Headers, lines...)
}

func (db *DB) Append(s section.Section, label string, fields ...any) {
	r := Record{Section: s, Label: label, Fields: fields}
	db.Records = append(db.Records, r)
}

func (db *DB) AppendRecors(records ...Record) {
	db.Records = append(db.Records, records...)
}
