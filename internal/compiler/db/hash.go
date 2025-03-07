package db

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func (db *DB) Hash() []byte {
	hasher := sha256.New()

	hasher.Write([]byte(strings.Join(db.headers, "")))

	for _, r := range db.regs {
		hasher.Write([]byte(fmt.Sprintf("%d", r.Section)))
		for _, f := range r.Fields {
			hasher.Write([]byte(f))
		}
	}

	return hasher.Sum(nil)
}
