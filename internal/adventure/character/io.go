package character

import (
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/vars"
	"thenewquill/internal/adventure/words"
	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"
)

func (s Store) Export(d *db.DB) {
	s.mut.Lock()
	defer s.mut.Unlock()

	for _, c := range s.chars {
		locationLabel := ""
		if c.Location != nil {
			locationLabel = c.Location.Label
		}

		d.Append(section.Chars, c.Label,
			c.Name.Label,
			c.Adjective.Label,
			c.Description,
			locationLabel,
			c.Created,
			c.Human,
			c.Vars.GetAll(),
		)
	}
}

func (s *Store) Import(d *db.DB, sw words.Store, locs loc.Store) {
	it := d.NewIterator(section.Chars)

	for {
		r := it.Next()
		if r == nil {
			break
		}

		c := &Character{
			Label:       r.Label,
			Name:        sw.Get(words.Noun, r.FieldAsString(0)),
			Adjective:   sw.Get(words.Adjective, r.FieldAsString(1)),
			Description: r.FieldAsString(2),
			Location:    locs.Get(r.FieldAsString(3)),
			Created:     r.FieldAsBool(4),
			Human:       r.FieldAsBool(5),
			Vars:        vars.NewStoreFromMap(r.FieldAsMap(6)),
		}

		s.Set(c)
	}
}
