package vars

import (
	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"
)

func (s *Store) Export(d *db.DB) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for k, v := range s.Regs {
		d.Append(section.Vars, k, v)
	}
}

func (s *Store) Import(d *db.DB) error {
	it := d.NewIterator(section.Vars)

	for {
		r := it.Next()
		if r == nil {
			break
		}

		s.Set(r.Label, r.Fields[0])
	}

	return nil
}
