package db

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha512"
	"encoding/gob"
	"errors"
	"io"
	"strings"
)

func (db *DB) Load(r io.Reader) error {
	db.Reset()

	br := bufio.NewReader(r)

	// headers
	for {
		b, err := br.Peek(1)
		if err != nil {
			if err == io.EOF {
				return ErrUnexpectedEOF
			}

			return errors.Join(ErrUnexpectedReadError, err)
		}

		if b[0] == EOH {
			_, err := br.Discard(1)
			if err != nil {
				return errors.Join(ErrUnexpectedReadError, err)
			}

			break
		}

		h, err := br.ReadString('\n')
		if err != nil {
			return err
		}

		db.AddHeader(strings.TrimSuffix(h, "\n"))
	}
	// crunched part
	zr, err := zlib.NewReader(br)
	if err != nil {
		return err
	}

	// records
	dec := gob.NewDecoder(zr)
	if err := dec.Decode(&db.Records); err != nil {
		return err
	}

	// end of records
	b, err := br.ReadByte()
	if err != nil {
		return errors.Join(ErrUnexpectedReadError, err)
	}

	if b != EOR {
		return ErrExpectedEOR
	}

	// hash
	hash := make([]byte, sha512.Size)
	if _, err := br.Read(hash); err != nil {
		return errors.Join(ErrUnexpectedReadError, err)
	}

	currentHash, err := db.Hash()
	if err != nil {
		return err
	}

	if !bytes.Equal(hash, currentHash) {
		return ErrHashMismatch
	}

	return nil
}
