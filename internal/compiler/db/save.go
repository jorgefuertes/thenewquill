package db

import (
	"compress/zlib"
	"encoding/gob"
	"io"
)

func (db *DB) Save(w io.Writer) error {
	// write headers
	for _, h := range db.Headers {
		if _, err := w.Write([]byte(h + "\n")); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte{EOH}); err != nil {
		return err
	}

	// crunched part
	zw := zlib.NewWriter(w)

	// write records
	enc := gob.NewEncoder(zw)
	if err := enc.Encode(db.Records); err != nil {
		return err
	}

	// flush crunched part
	if err := zw.Flush(); err != nil {
		return err
	}

	// end of records
	if _, err := w.Write([]byte{EOR}); err != nil {
		return err
	}

	// hash
	hash, err := db.Hash()
	if err != nil {
		return err
	}

	_, err = w.Write(hash)
	if err != nil {
		return err
	}

	return nil
}
