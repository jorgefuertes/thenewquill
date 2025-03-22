package words

import (
	"strings"

	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"
)

func (s Store) Export(d *db.DB) {
	for _, w := range s {
		d.Append(section.Words, w.Label,
			w.Type.Int(),
			strings.Join(w.Synonyms, ","),
		)
	}
}

func (s *Store) Import(d *db.DB) {
	it := d.NewIterator(section.Words)

	for {
		r := it.Next()
		if r == nil {
			break
		}

		label := r.Label
		t := WordType(r.FieldAsInt(0))
		syns := strings.Split(r.FieldAsString(1), ",")

		s.Set(label, t, syns...)
	}
}
