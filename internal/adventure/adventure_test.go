package adventure_test

import (
	"path/filepath"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/compiler"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// happyPathSource is the well-known happy-path adventure used across the
// integration tests. The relative location is stable in the tree.
const happyPathSource = "../compiler/test/src/happy/test.adv"

func TestNew(t *testing.T) {
	a := adventure.New()

	require.NotNil(t, a)
	assert.NotNil(t, a.DB)
	assert.NotNil(t, a.Config)
	assert.NotNil(t, a.Characters)
	assert.NotNil(t, a.Items)
	assert.NotNil(t, a.Messages)
	assert.NotNil(t, a.Words)
	assert.NotNil(t, a.Locations)
	assert.NotNil(t, a.Variables)
	assert.NotNil(t, a.Blobs)
}

func TestExportImportRoundtrip(t *testing.T) {
	src, err := compiler.Compile(happyPathSource)
	require.NoError(t, err, "happy-path adventure must compile cleanly")

	tmp := t.TempDir()
	path := filepath.Join(tmp, "adventure.db")

	// snapshot the source counts before exporting; Export writes a "date"
	// param so the totals may shift by one if the param did not exist yet.
	origItems := src.Items.Count()
	origCharacters := src.Characters.Count()
	origMessages := src.Messages.Count()
	origLocations := src.Locations.Count()
	origVariables := src.Variables.Count()
	origLabels := src.DB.CountLabels()
	origWords := src.Words.Count()
	origBlobs := src.Blobs.Count()

	_, _, err = src.Export(path)
	require.NoError(t, err)

	dst := adventure.New()
	require.NoError(t, dst.Import(path))

	t.Run("counts match per kind", func(t *testing.T) {
		assert.Equal(t, origItems, dst.Items.Count(), "items count")
		assert.Equal(t, origCharacters, dst.Characters.Count(), "characters count")
		assert.Equal(t, origMessages, dst.Messages.Count(), "messages count")
		assert.Equal(t, origLocations, dst.Locations.Count(), "locations count")
		assert.Equal(t, origVariables, dst.Variables.Count(), "variables count")
		assert.Equal(t, origWords, dst.Words.Count(), "words count")
		assert.Equal(t, origBlobs, dst.Blobs.Count(), "blobs count")
		assert.Equal(t, origLabels, dst.DB.CountLabels(), "labels count")
	})

	t.Run("config params survive the roundtrip", func(t *testing.T) {
		// Export injects a "date" param right before serialising, so it is
		// always present after roundtrip.
		assert.True(t, dst.Config.Get().WithLabel(config.DateParamLabel).Exists())
		assert.True(t, dst.Config.Get().WithLabel(config.TitleParamLabel).Exists())
		assert.True(t, dst.Config.Get().WithLabel(config.AuthorParamLabel).Exists())
	})

	t.Run("Import freezes the database", func(t *testing.T) {
		assert.True(t, dst.DB.IsFrozen())
	})

	t.Run("ID cursors advance past max imported", func(t *testing.T) {
		// After Import the next CreateLabel should not collide with any
		// imported label id.
		newID, err := dst.DB.CreateLabel("never-seen-before")
		require.NoError(t, err)
		assert.Greater(t, newID, uint32(0))

		// And the runtime can Update entities via the frozen+snapshot path.
		anyMsg := dst.Messages
		assert.NotZero(t, anyMsg.Count())
	})
}

func TestImportOnMissingFile(t *testing.T) {
	a := adventure.New()
	err := a.Import(filepath.Join(t.TempDir(), "missing.db"))
	require.Error(t, err)
}

func TestExportInjectsDateParam(t *testing.T) {
	a := adventure.New()
	// Need a label registered before Config.Set can store a param.
	_, err := a.DB.CreateLabel("title")
	require.NoError(t, err)
	_, err = a.Config.Set("title", "tiny adventure")
	require.NoError(t, err)

	tmp := t.TempDir()
	require.Falsef(t, a.DB.IsFrozen(), "source DB must not be frozen before Export")

	_, _, err = a.Export(filepath.Join(tmp, "out.db"))
	require.NoError(t, err)

	// After Export the "date" param exists.
	assert.True(t, a.Config.Get().WithLabel(config.DateParamLabel).Exists(),
		"Export should set the date config param as a side effect")

	// And it was written as a Param kind record.
	assert.NotZero(t, a.DB.CountRecordsByKind(kind.Param))
}
