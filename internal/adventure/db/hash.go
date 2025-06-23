package db

import (
	"crypto/sha256"
	"encoding/gob"
	"slices"
)

func (d *DB) Hash() ([]byte, error) {
	hasher := sha256.New()

	for _, l := range d.Labels {
		if err := gob.NewEncoder(hasher).Encode(l); err != nil {
			return []byte{}, err
		}
	}

	for _, r := range d.Data {
		if err := gob.NewEncoder(hasher).Encode(r); err != nil {
			return []byte{}, err
		}
	}

	return hasher.Sum(nil), nil
}

func (d *DB) IsHash(hash []byte) bool {
	h, err := d.Hash()
	if err != nil {
		return false
	}

	return slices.Equal(h, hash)
}
