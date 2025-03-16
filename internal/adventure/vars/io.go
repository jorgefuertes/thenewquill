package vars

import (
	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"
	"thenewquill/internal/util"
)

func (s *Store) Export(d *db.DB) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for k, v := range s.Regs {
		d.Add(db.NewRegister(section.Vars, k, util.ValueToString(v)))
	}
}

func (s *Store) Import(d *db.DB) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, r := range d.GetRegsForSection(section.Vars) {
		s.SetFromString(r.GetString(), r.GetString())
	}

	return nil
}
