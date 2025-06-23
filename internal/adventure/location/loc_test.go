package location_test

import (
	"errors"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocations(t *testing.T) {
	createWords := func(d *db.DB) {
		words := []word.Word{
			{Type: word.Noun, Synonyms: []string{"north", "n"}},
			{Type: word.Noun, Synonyms: []string{"east", "e"}},
			{Type: word.Noun, Synonyms: []string{"west", "w"}},
			{Type: word.Noun, Synonyms: []string{"south", "s"}},
		}

		for _, w := range words {
			id, err := d.Create(w.Synonyms[0], w)
			require.NoError(t, err, "cannot create word %q", w.Synonyms[0])
			assert.NotEmpty(t, id, "word %q should have an id", w.Synonyms[0])
		}
	}

	createLocations := func(d *db.DB) {
		definitions := []struct {
			label string
			title string
			desc  string
			conns map[string]string
		}{
			{"loc-001", "loc 001 title", "loc 001 desc", map[string]string{"east": "loc-002"}},
			{"loc-002", "loc 002 title", "loc 002 desc", map[string]string{"west": "loc-001"}},
			{"loc-003", "loc 003 title", "loc 003 desc", map[string]string{"north": "loc-004", "south": "loc-005"}},
			{"loc-004", "loc 004 title", "loc 004 desc", map[string]string{"east": "loc-003"}},
			{"loc-005", "loc 005 title", "loc 005 desc", map[string]string{"west": "loc-003"}},
			{"loc-006", "loc 006 title", "loc 006 desc", map[string]string{"north": "loc-007", "south": "loc-008"}},
			{"loc-007", "loc 007 title", "loc 007 desc", map[string]string{"east": "loc-006"}},
			{"loc-008", "loc 008 title", "loc 008 desc", map[string]string{"west": "loc-006"}},
			{"loc-009", "loc 009 title", "loc 009 desc", map[string]string{"north": "loc-001", "south": "loc-002"}},
		}

		for _, locDef := range definitions {
			l := location.New(locDef.title, locDef.desc)
			for k, v := range locDef.conns {
				var w word.Word
				err := d.GetByLabel(k, db.Words, &w)
				require.NoError(t, err)

				destLabel, err := d.GetLabelByName(v)
				if errors.Is(err, db.ErrNotFound) {
					// create the label
					var addErr error
					destLabel, addErr = d.AddLabel(v, false)
					require.NoError(t, addErr)
				} else if err != nil {
					require.NoError(t, err)
				}

				l.SetConn(w.GetID(), destLabel.ID)
			}

			_, err := d.Create(locDef.label, l)
			require.NoError(t, err)
		}
	}

	newDatabase := func() *db.DB {
		d := db.New()
		createWords(d)
		createLocations(d)

		return d
	}

	t.Run("create locations", func(t *testing.T) {
		d := newDatabase()
		require.NotZero(t, d.Count())
		require.NotZero(t, d.Query(db.Words).Count())
		require.NotZero(t, d.Query(db.Locations).Count())
	})
}
