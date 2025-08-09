package bin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/bin"
	"github.com/stretchr/testify/require"
)

func TestBinDB(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "test.adv")
	b := bin.New()

	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	t.Run("write label", func(t *testing.T) {
		l := db.Label{ID: 1, Name: "human"}

		err := b.WriteStoreable(l)
		require.NoError(t, err)
	})

	t.Run("write character", func(t *testing.T) {
		c := character.Character{
			ID:          1,
			NounID:      2,
			AdjectiveID: 3,
			Description: "description",
			LocationID:  4,
			Created:     false,
			Human:       true,
		}

		err := b.WriteStoreable(c)
		require.NoError(t, err)
	})

	t.Run("save to disk", func(t *testing.T) {
		err := b.Save(filename)
		require.NoError(t, err)
	})

	t.Run("load from disk", func(t *testing.T) {
		err := b.Load(filename)
		require.NoError(t, err)
	})

	t.Run("read label", func(t *testing.T) {
		l, err := b.ReadStoreable()
		require.NoError(t, err)
		require.IsType(t, &db.Label{}, l)
		require.Equal(t, db.ID(1), l.GetID())
	})

	t.Run("read character", func(t *testing.T) {
		s, err := b.ReadStoreable()
		require.NoError(t, err)
		require.IsType(t, &character.Character{}, s)
		require.Equal(t, db.ID(1), s.GetID())

		c, ok := s.(*character.Character)
		require.True(t, ok)
		require.Equal(t, db.ID(2), c.NounID)
		require.Equal(t, db.ID(3), c.AdjectiveID)
		require.Equal(t, "description", c.Description)
		require.Equal(t, db.ID(4), c.LocationID)
		require.Equal(t, false, c.Created)
		require.Equal(t, true, c.Human)
	})
}
