package words

import (
	"strings"

	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"
)

func (w Word) export() db.Register {
	return db.NewRegister(section.Words, w.Label,
		w.Label,
		w.Type.Int(),
		strings.Join(w.Synonyms, ","),
	)
}

func (s Store) Export(d *db.DB) {
	for _, w := range s {
		d.Add(w.export())
	}
}

func (s *Store) Import(d *db.DB) {
	for _, r := range d.GetRegsForSection(section.Words) {
		label := r.GetString()
		t := WordType(r.GetInt())
		syns := strings.Split(r.GetString(), ",")

		s.Set(label, t, syns...)
	}
}
