package bin

import (
	"encoding/gob"
	"fmt"
	"io"

	"thenewquill/internal/adventure"
	"thenewquill/internal/compiler/section"
)

func Export(a *adventure.Adventure, w io.Writer) error {
	dto := &adventureDTO{
		Config:    a.Config,
		Labels:    make(labels, 0),
		Vars:      a.Vars.GetAll(),
		Msgs:      make([]msgDTO, 0),
		Words:     make([]wordDTO, 0),
		Locations: make([]locationDTO, 0),
	}

	// messages
	for _, m := range a.Messages {
		id := dto.Labels.Add(m.Label, m.Type.Section())
		dto.Msgs = append(dto.Msgs, msgDTO{
			ID:      id,
			Text:    m.Text,
			Plurals: m.Plurals,
		})
	}

	// words
	for _, w := range a.Words {
		id := dto.Labels.Add(composeWordName(w.Label, w.Type), section.Words)
		dto.Words = append(dto.Words, wordDTO{
			ID:       id,
			Type:     w.Type,
			Synonyms: w.Synonyms,
		})
	}

	// create all the location labels first
	for _, l := range a.Locations {
		_ = dto.Labels.Add(l.Label, section.Locs)
	}

	// create the location DTOs in a second pass
	for _, l := range a.Locations {
		id := dto.Labels.Add(l.Label, section.Locs)
		d := locationDTO{ID: id, Title: l.Title, Description: l.Description, Conns: make(map[int]int)}
		for _, c := range l.Conns {
			wID, ok := dto.Labels.GetID(composeWordName(c.Word.Label, c.Word.Type), section.Words)
			if !ok {
				return fmt.Errorf("word label ID not found for %s", c.Word.Label+"#"+c.Word.Type.String())
			}

			destID, ok := dto.Labels.GetID(c.To.Label, section.Locs)
			if !ok {
				return fmt.Errorf("location label ID not found for %s", c.To.Label)
			}

			d.Conns[wID] = destID
		}

		dto.Locations = append(dto.Locations, d)
	}

	// chars
	for _, c := range a.Chars {
		id := dto.Labels.Add(c.Label, section.Chars)

		nameID, ok := dto.Labels.GetID(composeWordName(c.Name.Label, c.Name.Type), section.Words)
		if !ok {
			return fmt.Errorf("word label ID not found for %s", composeWordName(c.Name.Label, c.Name.Type))
		}

		adjectiveID, ok := dto.Labels.GetID(composeWordName(c.Adjective.Label, c.Adjective.Type), section.Words)
		if !ok {
			return fmt.Errorf("word label ID not found for %s", composeWordName(c.Adjective.Label, c.Adjective.Type))
		}

		locationID, ok := dto.Labels.GetID(c.Location.Label, section.Locs)
		if !ok {
			return fmt.Errorf("location label ID not found for %s", c.Location.Label)
		}

		dto.Chars = append(dto.Chars, charDTO{
			ID:          id,
			NameID:      nameID,
			AdjectiveID: adjectiveID,
			Description: c.Description,
			LocationID:  locationID,
			Created:     c.Created,
			Human:       c.Human,
		})
	}

	return gob.NewEncoder(w).Encode(dto)
}
