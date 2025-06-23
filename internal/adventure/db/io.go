package db

import (
	"encoding/gob"
	"io"
)

func (d *DB) Export(w io.Writer) error {
	return gob.NewEncoder(w).Encode(d)
}

func (d *DB) Import(r io.Reader) error {
	return gob.NewDecoder(r).Decode(d)
}
