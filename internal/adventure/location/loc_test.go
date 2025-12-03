package location_test

import (
	"fmt"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/database"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocations(t *testing.T) {
	db := database.New()
	wordStore := word.NewService(db)
	locationStore := location.NewService(db)

	t.Run("Words", func(t *testing.T) {
		wordsToCreate := [][]string{
			{"north", "n"},
			{"east", "e"},
			{"west", "w"},
			{"south", "s"},
		}

		// create words
		for _, syns := range wordsToCreate {
			labelID, err := db.CreateLabelIfNotExists(syns[0], false)
			require.NoError(t, err)

			w := word.New(word.Noun, syns...)
			w.SetLabelID(labelID)

			id, err := wordStore.Create(w)
			require.NoError(t, err, "cannot create word %q", syns[0])
			assert.False(t, id.IsUndefinedID(), "word %q should have an id", syns[0])
		}

		require.Equal(t, len(wordsToCreate), wordStore.Count())
	})

	t.Run("Locations", func(t *testing.T) {
		const (
			loc1 primitive.Label = "loc-001"
			loc2 primitive.Label = "loc-002"
			loc3 primitive.Label = "loc-003"
			loc4 primitive.Label = "loc-004"
			loc5 primitive.Label = "loc-005"
			loc6 primitive.Label = "loc-006"
			loc7 primitive.Label = "loc-007"
			loc8 primitive.Label = "loc-008"
			loc9 primitive.Label = "loc-009"
		)

		locData := []struct {
			label primitive.Label
			conns map[primitive.Label]primitive.Label
		}{
			{loc1, map[primitive.Label]primitive.Label{"east": loc2}},
			{loc2, map[primitive.Label]primitive.Label{"west": loc1}},
			{loc3, map[primitive.Label]primitive.Label{"north": loc4, "south": loc5}},
			{loc4, map[primitive.Label]primitive.Label{"east": loc3}},
			{loc5, map[primitive.Label]primitive.Label{"west": loc3}},
			{loc6, map[primitive.Label]primitive.Label{"north": loc7, "south": loc8}},
			{loc7, map[primitive.Label]primitive.Label{"east": loc6}},
			{loc8, map[primitive.Label]primitive.Label{"west": loc6}},
			{loc9, map[primitive.Label]primitive.Label{"north": loc1, "south": loc2}},
		}

		// create locations without connections
		for _, cur := range locData {
			title := fmt.Sprintf("%s title", cur.label)
			desc := fmt.Sprintf("%s desc", cur.label)
			loc := location.New(title, desc)

			labelID, err := db.CreateLabelIfNotExists(cur.label, false)
			require.NoError(t, err, "cannot create label %q", cur.label)

			loc.SetLabelID(labelID)

			id, err := locationStore.Create(loc)
			require.NoError(t, err, "creating location %q", cur.label)
			assert.True(t, id.IsDefinedID(), "location %q should have an id", cur.label)
		}

		require.Equal(t, len(locData), locationStore.Count())

		// create connections
		for _, cur := range locData {
			loc, err := locationStore.GetByLabel(cur.label)
			require.NoError(t, err, "cannot get location %q", cur.label)

			for wordLabel, dstLabel := range cur.conns {
				word, err := wordStore.GetByLabel(wordLabel)
				require.NoError(t, err, "cannot get destination word %q", wordLabel)

				dst, err := locationStore.GetByLabel(dstLabel)
				require.NoError(t, err, "cannot get destination location %q", dstLabel)

				loc.SetConn(word.ID, dst.ID)
				require.NoError(t, err, "cannot set connection %q->%q", wordLabel, dstLabel)
			}

			// update location
			err = locationStore.Update(loc)
			require.NoError(t, err, "cannot update location %q", cur.label)
		}
	})
}
