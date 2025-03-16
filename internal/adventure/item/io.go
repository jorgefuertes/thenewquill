package item

import (
	"fmt"
	"strings"

	"thenewquill/internal/adventure/character"
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/words"
	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"
	"thenewquill/internal/util"
)

func (s Store) Export(d *db.DB) {
	for _, item := range s.items {
		d.Add(item.export())
	}
}

func (i Item) export() db.Register {
	r := db.NewRegister(section.Items, i.Label,
		util.ToLabel(i.Noun),
		util.ToLabel(i.Adjective),
		i.Description,
		util.ValueToString(i.Weight),
		util.ValueToString(i.MaxWeight),
		util.ValueToString(i.IsContainer),
		util.ValueToString(i.IsWearable),
		util.ValueToString(i.IsCreated),
		util.ToLabel(i.Location),
		util.ToLabel(i.Inside),
		util.ToLabel(i.CarriedBy),
		util.ToLabel(i.WornBy),
	)

	// add vars at the end
	for k, v := range i.Vars.GetAll() {
		r.Fields = append(r.Fields, fmt.Sprintf("%s=%s", k, util.ValueToString(v)))
	}

	return r
}

func (s *Store) Import(d *db.DB, sw words.Store, locs loc.Store, cs character.Store) error {
	for _, r := range d.GetRegsForSection(section.Items) {
		label := r.GetString()
		noun := sw.Get(words.Noun, r.GetString())
		adj := sw.Get(words.Adjective, r.GetString())

		i := New(label, noun, adj)

		i.Description = r.GetString()
		i.Weight = r.GetInt()
		i.MaxWeight = r.GetInt()
		i.IsContainer = r.GetBool()
		i.IsWearable = r.GetBool()
		i.IsCreated = r.GetBool()
		i.Location = locs.Get(r.GetString())
		insideLabel := r.GetString()
		i.Inside = s.Get(insideLabel)
		i.CarriedBy = cs.Get(r.GetString())
		i.WornBy = cs.Get(r.GetString())

		for {
			v := r.GetString()
			if v == "" {
				break
			}

			parts := strings.Split(v, "=")
			if len(parts) != 2 {
				continue
			}

			i.Vars.Set(parts[0], parts[1])
		}

		if err := s.Set(i); err != nil {
			return err
		}
	}

	return nil
}
