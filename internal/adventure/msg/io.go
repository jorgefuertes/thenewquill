package msg

import (
	"github.com/jorgefuertes/thenewquill/internal/compiler/db"
	"github.com/jorgefuertes/thenewquill/internal/compiler/section"
)

func (s Store) Export(d *db.DB) {
	for _, m := range s {
		d.Append(section.Messages, m.Label,
			m.Text,
			m.Plurals[0],
			m.Plurals[1],
			m.Plurals[2],
		)
	}
}

func (s *Store) Import(d *db.DB) error {
	it := d.NewIterator(section.Messages)

	for {
		r := it.Next()
		if r == nil {
			break
		}

		err := s.Set(&Msg{
			Label: r.Label,
			Text:  r.FieldAsString(0),
			Plurals: [3]string{
				r.FieldAsString(1),
				r.FieldAsString(2),
				r.FieldAsString(3),
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
