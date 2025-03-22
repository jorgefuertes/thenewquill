package loc

import (
	"fmt"

	"thenewquill/internal/adventure/words"
	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"
)

func (s Store) Export(d *db.DB) {
	for _, l := range s {
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

		l := s.Set(r.Label, r.FieldAsString(0), r.FieldAsString(1))
		l.Vars.SetAll(r.FieldAsMap(3))

		for k, v := range r.FieldAsMap(2) {
			w := sw.First(k)
			if w == nil {
				return fmt.Errorf("cannot find word %s for connection in location %s", k, l.Label)
			}

			toLabel, ok := v.(string)
			if !ok {
				return fmt.Errorf("cannot convert value %v to string for connection in location %s", v, l.Label)
			}

			to := s.Get(toLabel)
			if to == nil {
				to = s.Set(toLabel, Undefined, Undefined)
			}

			l.SetConn(w, to)
		}
	}

	return nil
}
