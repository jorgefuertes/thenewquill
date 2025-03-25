package loc

import (
	"fmt"

	"thenewquill/internal/adventure/words"
	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"
)

func (s Store) Export(d *db.DB) {
	s.mut.Lock()
	defer s.mut.Unlock()

	for _, l := range s.data {
		cMap := make(map[string]string, 0)
		for _, c := range l.Conns {
			cMap[c.Word.Label] = c.To.Label
		}

		d.Append(section.Locs, l.Label,
			l.Title,
			l.Description,
			cMap,
			l.Vars.GetAll(),
		)
	}
}

func (s *Store) Import(d *db.DB, sw words.Store) error {
	it := d.NewIterator(section.Locs)

	for {
		r := it.Next()
		if r == nil {
			break
		}

		l := New(r.Label, r.FieldAsString(0), r.FieldAsString(1))
		l.Vars.SetAll(r.FieldAsMapAny(3))

		for wordLabel, toLabel := range r.FieldAsMapString(2) {
			w := sw.First(wordLabel)
			if w == nil {
				return fmt.Errorf("cannot find word %s for connection in location %s", wordLabel, l.Label)
			}

			to := s.Get(toLabel)
			if to == nil {
				var err error

				if to, err = s.New(toLabel); err != nil {
					return err
				}
			}

			l.SetConn(w, to)
		}

		if err := s.Set(l); err != nil {
			return err
		}
	}

	return nil
}
