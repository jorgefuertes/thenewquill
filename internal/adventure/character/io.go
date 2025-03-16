package character

import (
	"fmt"
	"strings"

	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/vars"
	"thenewquill/internal/adventure/words"
	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"
	"thenewquill/internal/util"
)

func (s Store) Export(d *db.DB) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, c := range s.chars {
		d.Add(c.export())
	}
}

func (c Character) export() db.Register {
	locationLabel := ""
	if c.Location != nil {
		locationLabel = c.Location.Label
	}

	r := db.NewRegister(section.Chars, c.Label,
		c.Name.Label,
		c.Adjective.Label,
		c.Description,
		locationLabel,
		util.ValueToString(c.Created),
		util.ValueToString(c.Human),
	)

	for k, v := range c.Vars.GetAll() {
		r.Fields = append(r.Fields, fmt.Sprintf("%s=%s", k, util.ValueToString(v)))
	}

	return r
}

func (s *Store) Import(d *db.DB, sw words.Store, locs loc.Store) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, r := range d.GetRegsForSection(section.Chars) {
		c := &Character{
			Label:       r.GetString(),
			Name:        sw.Get(words.Noun, r.GetString()),
			Adjective:   sw.Get(words.Adjective, r.GetString()),
			Description: r.GetString(),
			Location:    locs.Get(r.GetString()),
			Created:     r.GetBool(),
			Human:       r.GetBool(),
			Vars:        vars.NewStore(),
		}

		for {
			v := r.GetString()
			if v == "" {
				break
			}

			parts := strings.Split(v, "=")
			if len(parts) != 2 {
				continue
			}

			c.Vars.Set(parts[0], parts[1])
		}

		if err := s.Set(c); err != nil {
			return err
		}
	}

	return nil
}
