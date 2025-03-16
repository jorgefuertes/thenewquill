package loc

import (
	"strings"

	"thenewquill/internal/adventure/words"
	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"
)

func (l Location) export() db.Register {
	conns := make([]string, 0)
	for _, c := range l.Conns {
		conns = append(conns, c.Word.Label+":"+c.To.Label)
	}

	return db.NewRegister(section.Locs, l.Label,
		l.Title,
		l.Description,
		conns,
	)
}

func (s Store) Export(d *db.DB) {
	for _, l := range s {
		d.Add(l.export())
	}
}

func (s *Store) Import(d *db.DB, sw words.Store) {
	for _, r := range d.GetRegsForSection(section.Locs) {
		l := s.Set(r.GetString(), r.GetString(), r.GetString())

		conns := strings.Split(r.GetString(), ",")
		for _, conn := range conns {
			parts := strings.Split(conn, ":")
			if len(parts) != 2 {
				continue
			}

			w := sw.First(parts[0])
			to := s.Get(parts[1])
			if to == nil {
				to = s.Set(parts[1], "", "")
			}

			l.SetConn(w, to)
		}
	}
}
