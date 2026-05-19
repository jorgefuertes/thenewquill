package database

import (
	"path/filepath"
	"strconv"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// seedTestItem creates one persisted testItem with the given label and
// returns it.
func seedTestItem(t *testing.T, db *DB, label string) *testItem {
	t.Helper()

	labelID, err := db.CreateLabel(label)
	require.NoError(t, err)

	it := &testItem{
		LabelID: labelID,
		Title:   "Title for " + label,
		Weight:  10,
	}

	_, err = db.Create(it)
	require.NoError(t, err)

	return it
}

func TestCreateEdgeCases(t *testing.T) {
	t.Run("frozen DB rejects Create", func(t *testing.T) {
		db := NewDB()
		db.Freeze()

		labelID, err := db.CreateLabel("foo")
		require.NoError(t, err)

		_, err = db.Create(&testItem{LabelID: labelID})
		require.ErrorIs(t, err, ErrDatabaseIsFrozen)
	})

	t.Run("non-zero ID is rejected", func(t *testing.T) {
		db := NewDB()
		labelID, _ := db.CreateLabel("foo")

		_, err := db.Create(&testItem{ID: 42, LabelID: labelID})
		require.ErrorIs(t, err, ErrIDFieldIsNotZero)
	})

	t.Run("unknown LabelID is rejected", func(t *testing.T) {
		db := NewDB()
		_, err := db.Create(&testItem{LabelID: 99999})
		require.ErrorIs(t, err, ErrLabelNotFound)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("missing ID", func(t *testing.T) {
		db := NewDB()
		err := db.Update(&testItem{LabelID: 1})
		require.ErrorIs(t, err, ErrMissingIDToUpdate)
	})

	t.Run("unknown LabelID", func(t *testing.T) {
		db := NewDB()
		it := seedTestItem(t, db, "x")
		it.LabelID = 99999

		err := db.Update(it)
		require.ErrorIs(t, err, ErrLabelNotFound)
	})

	t.Run("record not found", func(t *testing.T) {
		db := NewDB()
		labelID, _ := db.CreateLabel("foo")
		err := db.Update(&testItem{ID: 999, LabelID: labelID})
		require.ErrorIs(t, err, ErrRecordNotFound)
	})

	t.Run("happy path persists changes", func(t *testing.T) {
		db := NewDB()
		it := seedTestItem(t, db, "x")

		it.Title = "updated"
		require.NoError(t, db.Update(it))

		var got testItem
		require.NoError(t, db.Get(it.ID, &got))
		assert.Equal(t, "updated", got.Title)
	})
}

func TestCountRecordsByKind(t *testing.T) {
	db := NewDB()

	t.Run("empty db", func(t *testing.T) {
		assert.Zero(t, db.CountRecordsByKind(kind.Test))
	})

	t.Run("after Create", func(t *testing.T) {
		seedTestItem(t, db, "one")
		seedTestItem(t, db, "two")
		assert.Equal(t, 2, db.CountRecordsByKind(kind.Test))
	})

	t.Run("kind.Label is a special case that returns label count", func(t *testing.T) {
		// 3 reserved labels (#undefined, *, _) + 2 created above.
		assert.Equal(t, 5, db.CountRecordsByKind(kind.Label))
	})

	t.Run("unrelated kind stays zero", func(t *testing.T) {
		assert.Zero(t, db.CountRecordsByKind(kind.Item))
	})
}

func TestGetByLabelID(t *testing.T) {
	db := NewDB()
	it := seedTestItem(t, db, "sword")

	var got testItem
	require.NoError(t, db.GetByLabelID(it.LabelID, &got))
	assert.Equal(t, it.ID, got.ID)

	require.ErrorIs(t, db.GetByLabelID(99999, &got), ErrRecordNotFound)
}

func TestLabelHelpers(t *testing.T) {
	db := NewDB()

	t.Run("GetLabelOrBlank with valid id", func(t *testing.T) {
		assert.Equal(t, LabelAsterisk, db.GetLabelOrBlank(1))
		assert.Equal(t, LabelUnderscore, db.GetLabelOrBlank(2))
	})

	t.Run("GetLabelOrBlank with unknown id returns blank", func(t *testing.T) {
		assert.Empty(t, db.GetLabelOrBlank(99999))
	})

	t.Run("GetLabelFromRecordOrBlank traverses record -> label", func(t *testing.T) {
		it := seedTestItem(t, db, "key")
		assert.Equal(t, "key", db.GetLabelFromRecordOrBlank(it.ID))
	})

	t.Run("GetLabelFromRecordOrBlank with unknown record id returns blank", func(t *testing.T) {
		assert.Empty(t, db.GetLabelFromRecordOrBlank(99999))
	})
}

func TestNewFilter(t *testing.T) {
	f := NewFilter("Title", Contains, "hello")
	assert.Equal(t, "Title", f.field)
	assert.Equal(t, Contains, f.condition)
	assert.Equal(t, "hello", f.value)
}

func TestCursorClose(t *testing.T) {
	c := newCursor()
	c.addOrReplace(Record{LabelID: 1, Kind: kind.Test})
	require.Equal(t, 1, c.Count())

	c.Close()
	assert.Equal(t, 0, c.Count())
	assert.Nil(t, c.data, "Close should release the backing slice")
}

func TestFreezeAndSnapshot(t *testing.T) {
	t.Run("Snapshot is a no-op when not frozen", func(t *testing.T) {
		db := NewDB()
		db.Snapshot()
		assert.Empty(t, db.snapshots)
	})

	t.Run("SnapBack returns false when not frozen", func(t *testing.T) {
		db := NewDB()
		assert.False(t, db.SnapBack())
	})

	t.Run("SnapBack returns false when frozen but no snapshots", func(t *testing.T) {
		db := NewDB()
		db.Freeze()
		assert.False(t, db.SnapBack())
	})

	t.Run("Freeze + Snapshot + SnapBack lifecycle", func(t *testing.T) {
		db := NewDB()
		seedTestItem(t, db, "first")

		db.Freeze()
		assert.True(t, db.IsFrozen())

		db.Snapshot()
		require.Len(t, db.snapshots, 1)

		db.Snapshot()
		require.Len(t, db.snapshots, 2)

		assert.True(t, db.SnapBack())
		assert.Len(t, db.snapshots, 1)

		assert.True(t, db.SnapBack())
		assert.Empty(t, db.snapshots)

		assert.False(t, db.SnapBack(), "no more snapshots to pop")
	})
}

func TestUpdateInFrozenDB(t *testing.T) {
	t.Run("Update on frozen DB without snapshot auto-creates one", func(t *testing.T) {
		db := NewDB()
		it := seedTestItem(t, db, "x")
		db.Freeze()
		require.Empty(t, db.snapshots)

		it.Title = "updated"
		require.NoError(t, db.Update(it))

		require.Len(t, db.snapshots, 1, "Update must implicitly create a snapshot")
		require.Contains(t, db.snapshots[0], it.ID)
	})

	t.Run("Update writes to last snapshot when frozen", func(t *testing.T) {
		db := NewDB()
		it := seedTestItem(t, db, "x")
		db.Freeze()
		db.Snapshot()

		it.Title = "updated in snapshot"
		require.NoError(t, db.Update(it))

		// The original record is still in db.data; the updated copy lives in
		// the last snapshot.
		require.Contains(t, db.data, it.ID, "base data still holds the record")
		require.Contains(t, db.snapshots[len(db.snapshots)-1], it.ID,
			"snapshot must hold the updated record")
	})
}

func TestRamSaveAndLoad(t *testing.T) {
	db := NewDB()
	it := seedTestItem(t, db, "x")
	db.Freeze()
	db.Snapshot()

	it.Title = "after-update"
	require.NoError(t, db.Update(it))

	db.RamSave()
	require.NotEmpty(t, db.ram, "RamSave flattens all snapshots into ram")

	// RamLoad appends the ram snapshot as a new layer.
	before := len(db.snapshots)
	db.RamLoad()
	assert.Equal(t, before+1, len(db.snapshots),
		"RamLoad appends a snapshot rather than replacing")
}

func TestImportPreservesIDCursors(t *testing.T) {
	// Build a DB with several labels and records, export, then import into a
	// fresh DB. After import, the next CreateLabel and Create must yield IDs
	// strictly greater than any imported ID, so we cannot accidentally
	// overwrite an imported record/label.
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.db")

	src := NewDB()
	for i := range 5 {
		seedTestItem(t, src, "rec-"+strconv.Itoa(i))
	}

	maxLabelID := src.lastLabelID
	maxDataID := src.lastDataID

	_, _, err := src.Export(path)
	require.NoError(t, err)

	dst := NewDB()
	require.NoError(t, dst.Import(path))

	assert.Equal(t, maxLabelID, dst.lastLabelID, "lastLabelID must match max imported label id")
	assert.Equal(t, maxDataID, dst.lastDataID, "lastDataID must match max imported record id")

	// Next CreateLabel returns a fresh id, not colliding with any imported one.
	newLabelID, err := dst.CreateLabel("brand-new")
	require.NoError(t, err)
	assert.Greater(t, newLabelID, maxLabelID, "new label id must be greater than imported max")
	assert.NotContains(t, []uint32{0, 1, 2, 3, 4, 5}, newLabelID, "and definitely not overwriting an imported one")

	// And nothing was clobbered.
	for id, origLabel := range src.labels {
		got, ok := dst.labels[id]
		require.True(t, ok)
		assert.Equal(t, origLabel, got)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "save.dat")

	t.Run("Save with no snapshots fails", func(t *testing.T) {
		db := NewDB()
		_, _, err := db.Save(path)
		require.ErrorIs(t, err, ErrNothingToSave)
	})

	t.Run("Save and Load roundtrip", func(t *testing.T) {
		src := NewDB()
		seedTestItem(t, src, "x")
		src.Freeze()
		src.Snapshot()

		// Drop a record into the snapshot via Update so Save has something
		// meaningful to write.
		all := src.Query(FilterByKind(kind.Test))
		var first testItem
		require.NoError(t, all.First(&first))
		first.Title = "after save"
		require.NoError(t, src.Update(&first))

		_, _, err := src.Save(path)
		require.NoError(t, err)

		dst := NewDB()
		require.NoError(t, dst.Load(path))
		assert.Equal(t, len(src.snapshots), len(dst.snapshots))
	})

	t.Run("Load on a missing file returns error", func(t *testing.T) {
		dst := NewDB()
		require.Error(t, dst.Load(filepath.Join(tmp, "missing.dat")))
	})
}

func TestQueryFrozenSnapshotsOverlay(t *testing.T) {
	// After an Update on a frozen DB, Query must surface the snapshot copy
	// (not the original from base data), with no duplicates.
	db := NewDB()
	it := seedTestItem(t, db, "x")
	db.Freeze()
	db.Snapshot()

	it.Title = "v2"
	require.NoError(t, db.Update(it))

	c := db.Query(FilterByKind(kind.Test))
	require.Equal(t, 1, c.Count(), "snapshot overlay must not duplicate the base record")

	var got testItem
	require.NoError(t, c.First(&got))
	assert.Equal(t, "v2", got.Title, "snapshot overlay must win over base data")
}

func TestQueryMergesAcrossMultipleSnapshots(t *testing.T) {
	// Two snapshots updating the same record: the most recent one wins.
	db := NewDB()
	it := seedTestItem(t, db, "x")
	db.Freeze()

	db.Snapshot()
	it.Title = "v2"
	require.NoError(t, db.Update(it))

	db.Snapshot()
	it.Title = "v3"
	require.NoError(t, db.Update(it))

	c := db.Query(FilterByKind(kind.Test))
	require.Equal(t, 1, c.Count())

	var got testItem
	require.NoError(t, c.First(&got))
	assert.Equal(t, "v3", got.Title, "most recent snapshot wins")
}
