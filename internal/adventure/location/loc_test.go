package location_test

import (
	"fmt"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"github.com/jorgefuertes/thenewquill/internal/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocations(t *testing.T) {
	db := database.NewDB()
	wordStore := word.NewService(db)
	locationStore := location.NewService(db)
	prevCount := wordStore.Count()

	t.Run("Words", func(t *testing.T) {
		wordsToCreate := [][]string{
			{"north", "n"},
			{"east", "e"},
			{"west", "w"},
			{"south", "s"},
		}

		// create words
		for _, syns := range wordsToCreate {
			labelID, err := db.CreateLabel(syns[0])
			require.NoError(t, err)

			w := word.New(labelID, word.Noun, syns...)

			id, err := wordStore.Create(w)
			require.NoError(t, err, "cannot create word %q", syns[0])
			assert.NotZero(t, id, "word %q should have an id", syns[0])
		}

		require.Equal(t, len(wordsToCreate)+prevCount, wordStore.Count())
	})

	t.Run("Locations", func(t *testing.T) {
		const (
			loc1 = "loc-001"
			loc2 = "loc-002"
			loc3 = "loc-003"
			loc4 = "loc-004"
			loc5 = "loc-005"
			loc6 = "loc-006"
			loc7 = "loc-007"
			loc8 = "loc-008"
			loc9 = "loc-009"
		)

		locData := []struct {
			label string
			conns map[string]string
		}{
			{loc1, map[string]string{"east": loc2}},
			{loc2, map[string]string{"west": loc1}},
			{loc3, map[string]string{"north": loc4, "south": loc5}},
			{loc4, map[string]string{"east": loc3}},
			{loc5, map[string]string{"west": loc3}},
			{loc6, map[string]string{"north": loc7, "south": loc8}},
			{loc7, map[string]string{"east": loc6}},
			{loc8, map[string]string{"west": loc6}},
			{loc9, map[string]string{"north": loc1, "south": loc2}},
		}

		// create locations without connections
		for _, cur := range locData {
			loc := location.New()
			loc.Title = fmt.Sprintf("%s title", cur.label)
			loc.Description = fmt.Sprintf("%s desc", cur.label)

			labelID, err := db.CreateLabel(cur.label)
			require.NoError(t, err, "cannot create label %q", cur.label)

			loc.LabelID = labelID

			id, err := locationStore.Create(loc)
			require.NoError(t, err, "creating location %q", cur.label)
			assert.NotZero(t, id, "location %q should have an id", cur.label)
		}

		require.Equal(t, len(locData), locationStore.Count())

		// create connections
		for _, cur := range locData {
			loc, err := locationStore.Get().WithLabel(cur.label).First()
			require.NoError(t, err, "cannot get location %q", cur.label)

			for wordLabel, dstLabel := range cur.conns {
				word, err := wordStore.Get().WithLabel(wordLabel).First()
				require.NoError(t, err, "cannot get destination word %q", wordLabel)

				dst, err := locationStore.Get().WithLabel(dstLabel).First()
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
