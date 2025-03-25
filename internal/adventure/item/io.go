package item

import (
	"thenewquill/internal/adventure/character"
	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/vars"
	"thenewquill/internal/adventure/words"
	"thenewquill/internal/compiler/db"
	"thenewquill/internal/compiler/section"
)

func (s Store) Export(d *db.DB) {
	for _, item := range s.items {
		locationLabel := ""
		if item.Location != nil {
			locationLabel = item.Location.Label
		}

		insideLabel := ""
		if item.Inside != nil {
			insideLabel = item.Inside.Label
		}

		carriedByLabel := ""
		if item.CarriedBy != nil {
			carriedByLabel = item.CarriedBy.Label
		}

		wornByLabel := ""
		if item.WornBy != nil {
			wornByLabel = item.WornBy.Label
		}

		d.Append(section.Items, item.Label,
			item.Noun.Label,
			item.Adjective.Label,
			item.Description,
			item.Weight,
			item.MaxWeight,
			item.IsContainer,
			item.IsWearable,
			item.IsCreated,
			locationLabel,
			insideLabel,
			carriedByLabel,
			wornByLabel,
			item.Vars.GetAll(),
		)
	}
}

func (s *Store) Import(d *db.DB, sw words.Store, locs loc.Store, cs character.Store) error {
	it := d.NewIterator(section.Items)

	for {
		r := it.Next()
		if r == nil {
			break
		}

		location := locs.Get(r.FieldAsString(8))
		if location == nil && r.FieldAsString(8) != "" {
			var err error
			location, err = locs.New(r.FieldAsString(8))
			if err != nil {
				return err
			}
		}

		insideContainer := s.Get(r.FieldAsString(9))
		if insideContainer == nil && r.FieldAsString(9) != "" {
			var err error
			insideContainer, err = s.New(r.FieldAsString(9))
			if err != nil {
				return err
			}
		}

		carriedByChar := cs.Get(r.FieldAsString(10))
		if carriedByChar == nil {
			carriedByChar = cs.New(r.FieldAsString(10))
		}

		wornByChar := cs.Get(r.FieldAsString(11))
		if wornByChar == nil {
			wornByChar = cs.New(r.FieldAsString(11))
		}

		i := &Item{
			Label:       r.Label,
			Noun:        sw.Get(words.Noun, r.FieldAsString(0)),
			Adjective:   sw.Get(words.Adjective, r.FieldAsString(1)),
			Description: r.FieldAsString(2),
			Weight:      r.FieldAsInt(3),
			MaxWeight:   r.FieldAsInt(4),
			IsContainer: r.FieldAsBool(5),
			IsWearable:  r.FieldAsBool(6),
			IsCreated:   r.FieldAsBool(7),
			Location:    location,
			Inside:      insideContainer,
			CarriedBy:   carriedByChar,
			WornBy:      wornByChar,
			Vars:        vars.NewStoreFromMap(r.FieldAsMapAny(12)),
		}

		if err := s.Set(i); err != nil {
			return err
		}
	}

	return nil
}
