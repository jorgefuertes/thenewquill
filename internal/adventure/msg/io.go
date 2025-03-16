package msg

import (
	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"
)

func (m Msg) export() db.Register {
	return db.NewRegister(section.Messages, m.Label,
		m.Text,
		m.Plurals[0],
		m.Plurals[1],
		m.Plurals[2],
	)
}

func (s Store) Export(d *db.DB) {
	for _, m := range s {
		d.Add(m.export())
	}
}

func (s *Store) Import(d *db.DB) {
	for _, r := range d.GetRegsForSection(section.Messages) {
		_ = s.Set(&Msg{
			Label: r.GetString(),
			Text:  r.GetString(),
			Plurals: [3]string{
				r.GetString(),
				r.GetString(),
				r.GetString(),
			},
		})
	}
}
