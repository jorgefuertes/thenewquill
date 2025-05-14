package db_test

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockStoreable struct {
	id      db.ID
	kind    db.Kind
	subKind db.SubKind
	Name    string
}

func (m MockStoreable) GetID() db.ID                   { return m.id }
func (m MockStoreable) GetKind() (db.Kind, db.SubKind) { return m.kind, m.subKind }
func (m MockStoreable) Validate() error                { return nil }

const maxRecordsPerSubKind = 5

func generateDatabase(t *testing.T) *db.DB {
	t.Helper()

	database := db.New()

	for _, k := range db.Kinds() {
		if k == db.None {
			continue
		}

		for _, sk := range k.SubKinds() {
			for i := 0; i < maxRecordsPerSubKind; i++ {
				name := fmt.Sprintf("%s-%s-%d",
					strings.ReplaceAll(k.String(), " ", "_"),
					strings.ReplaceAll(sk.String(), " ", "_"),
					i)
				l, err := database.AddLabel(name)
				require.NoError(t, err)
				mock := MockStoreable{id: l.ID, kind: k, subKind: sk, Name: l.Name}
				require.NoError(t, database.Append(mock))
			}
		}
	}

	return database
}

func TestDB(t *testing.T) {
	t.Run("New DB initialization", func(t *testing.T) {
		database := db.New()
		assert.NotNil(t, database, "New() should not return nil")
		assert.Equal(t, 0, database.Len(), "New database should be empty")
	})

	t.Run("Get", func(t *testing.T) {
		database := generateDatabase(t)

		var firstVerbID db.ID

		t.Run("GetByKind", func(t *testing.T) {
			retrieved := database.GetByKind(db.Words, db.VerbSubKind)
			assert.Len(t, retrieved, maxRecordsPerSubKind)
			firstVerbID = retrieved[0].GetID()
		})

		t.Run("Get one by ID", func(t *testing.T) {
			retrieved, err := database.Get(db.Words, db.VerbSubKind, firstVerbID)
			require.NoError(t, err)
			assert.NotEmpty(t, retrieved)
		})

		t.Run("GetAs", func(t *testing.T) {
			var retrievedAs MockStoreable
			err := database.GetAs(firstVerbID, db.Words, db.VerbSubKind, &retrievedAs)
			require.NoError(t, err)
			assert.NotEmpty(t, retrievedAs)

			err = database.GetAs(firstVerbID, db.Words, db.VerbSubKind, retrievedAs)
			assert.ErrorIs(t, err, db.ErrDstMustBePointer)

			err = database.GetAs(firstVerbID, db.Words, db.VerbSubKind, nil)
			assert.ErrorIs(t, err, db.ErrDstMustBePointer)

			err = database.GetAs(65_535, db.Words, db.VerbSubKind, &retrievedAs)
			assert.ErrorIs(t, err, db.ErrRecordNotFound)
		})

		t.Run("can't append with subkind any", func(t *testing.T) {
			mock := MockStoreable{id: 5, kind: db.Words, subKind: db.AnySubKind}
			err := database.Append(mock)
			assert.ErrorIs(t, err, db.ErrSubKindMustBeDefined)
		})

		t.Run("can't append with kind none", func(t *testing.T) {
			mock := MockStoreable{id: 5, kind: db.None, subKind: db.NoSubKind}
			err := database.Append(mock)
			assert.ErrorIs(t, err, db.ErrKindCannotBeNone)
		})

		t.Run("can't get non existent record", func(t *testing.T) {
			_, err := database.Get(db.Words, db.VerbSubKind, 6)
			assert.ErrorIs(t, err, db.ErrRecordNotFound)
		})
	})

	t.Run("Duplicate records", func(t *testing.T) {
		database := db.New()
		mock := MockStoreable{id: 5, kind: db.Words, subKind: db.VerbSubKind}

		_ = database.Append(mock)
		err := database.Append(mock)
		assert.ErrorIs(t, err, db.ErrDuplicatedRecord)
	})

	t.Run("Update record", func(t *testing.T) {
		database := db.New()
		mock := MockStoreable{id: 5, kind: db.Words, subKind: db.VerbSubKind, Name: "test"}
		updatedMock := MockStoreable{id: 5, kind: db.Words, subKind: db.VerbSubKind, Name: "test2"}

		require.NoError(t, database.Append(mock))
		require.NoError(t, database.Update(updatedMock))

		retrieved, err := database.Get(db.Words, db.VerbSubKind, 5)
		require.NoError(t, err)

		k, sk := retrieved.GetKind()
		assert.Equal(t, db.Words, k)
		assert.Equal(t, db.VerbSubKind, sk)

		r := reflect.ValueOf(retrieved)
		assert.Equal(t, "test2", r.FieldByName("Name").String())

		t.Run("can't update non existent record", func(t *testing.T) {
			mock := MockStoreable{id: 6, kind: db.Words, subKind: db.VerbSubKind, Name: "test"}
			err := database.Update(mock)
			assert.ErrorIs(t, err, db.ErrRecordNotFound)
		})
	})

	t.Run("Remove record", func(t *testing.T) {
		database := db.New()
		mock := MockStoreable{id: 5, kind: db.Words, subKind: db.VerbSubKind}

		require.NoError(t, database.Append(mock))
		require.NoError(t, database.Remove(5, db.Words, db.VerbSubKind))
		assert.False(t, database.Exists(5, db.Words, db.VerbSubKind))
	})

	t.Run("GetByKind", func(t *testing.T) {
		database := generateDatabase(t)
		words := database.GetByKind(db.Words, db.VerbSubKind)
		assert.Len(t, words, 5)
	})

	t.Run("GetByKindAs", func(t *testing.T) {
		database := generateDatabase(t)

		var words []MockStoreable
		err := database.GetByKindAs(db.Words, db.VerbSubKind, words)

		require.Error(t, err)
		assert.ErrorIs(t, err, db.ErrDstMustBePointerSlice)

		err = database.GetByKindAs(db.Words, db.VerbSubKind, &words)
		require.NoError(t, err)
		assert.Len(t, words, 5)

		var noCasteable []string
		err = database.GetByKindAs(db.Words, db.VerbSubKind, &noCasteable)
		assert.ErrorIs(t, err, db.ErrCannotCastFromStoreable)
	})

	t.Run("Remove", func(t *testing.T) {
		database := generateDatabase(t)

		items := database.GetByKind(db.Words, db.VerbSubKind)
		require.NotEmpty(t, items)

		firstID := items[0].GetID()
		err := database.Remove(firstID, db.Words, db.VerbSubKind)
		require.NoError(t, err)
		assert.False(t, database.Exists(firstID, db.Words, db.VerbSubKind))

		err = database.Remove(firstID, db.Words, db.VerbSubKind)
		assert.ErrorIs(t, err, db.ErrRecordNotFound)
	})

	t.Run("Count and Reset", func(t *testing.T) {
		database := generateDatabase(t)

		assert.Equal(t, maxRecordsPerSubKind, database.Count(db.Words, db.VerbSubKind))
		database.Reset()
		assert.Equal(t, 0, database.Count(db.Words, db.VerbSubKind))
		assert.Equal(t, 0, database.CountAll())
	})
}

func TestLabels(t *testing.T) {
	t.Run("Default labels initialization", func(t *testing.T) {
		database := db.New()
		defaultLabels := []string{"undefined", "adv-config", "_", "*"}

		for _, name := range defaultLabels {
			assert.True(t, database.ExistsLabelName(name),
				"Default label %s should exist", name)
		}
	})

	t.Run("Add valid label", func(t *testing.T) {
		database := db.New()
		label, err := database.AddLabel("test-label")

		require.NoError(t, err)
		assert.Equal(t, "test-label", label.Name)
		assert.GreaterOrEqual(t, label.ID, db.MinMeaningfulID)
	})

	t.Run("Add invalid label", func(t *testing.T) {
		database := db.New()
		invalidNames := []string{
			"test label",
			"test!label",
			"@label",
			"",
		}

		for _, name := range invalidNames {
			_, err := database.AddLabel(name)
			assert.ErrorIs(t, err, db.ErrInvalidLabelName,
				"Should fail for invalid name: %s", name)
		}
	})

	t.Run("Get label by ID", func(t *testing.T) {
		database := db.New()
		added, err := database.AddLabel("test-label")
		require.NoError(t, err)

		retrieved, err := database.GetLabel(added.ID)
		require.NoError(t, err)
		assert.Equal(t, "test-label", retrieved.Name)
	})

	t.Run("Get non-existent label", func(t *testing.T) {
		database := db.New()
		_, err := database.GetLabel(999)
		assert.ErrorIs(t, err, db.ErrNotFound)
	})

	t.Run("Get label by name", func(t *testing.T) {
		database := db.New()
		added, err := database.AddLabel("test-label")
		require.NoError(t, err)

		retrieved, err := database.GetLabelByName("test-label")
		require.NoError(t, err)
		assert.Equal(t, added.ID, retrieved.ID)

		_, err = database.GetLabelByName("test-label-2")
		assert.ErrorIs(t, err, db.ErrNotFound)

		t.Run("Get label name", func(t *testing.T) {
			name := database.GetLabelName(added.ID)
			assert.Equal(t, "test-label", name)

			name = database.GetLabelName(math.MaxInt32)
			assert.Equal(t, db.UndefinedLabel.Name, name)
		})
	})

	t.Run("Duplicate label prevention", func(t *testing.T) {
		database := db.New()
		label1, err := database.AddLabel("test-label")
		require.NoError(t, err)

		label2, err := database.AddLabel("test-label")
		require.NoError(t, err)

		assert.Equal(t, label1.ID, label2.ID)
	})

	t.Run("Valid label names", func(t *testing.T) {
		validNames := []string{
			"simple",
			"test-label",
			"test_label",
			"label123",
			"123label",
			"étiquette",
		}

		for _, name := range validNames {
			assert.True(t, db.IsValidLabelName(name),
				"Should be valid label name: %s", name)
		}
	})

	t.Run("Label ID conversions", func(t *testing.T) {
		id := db.ID(42)
		assert.Equal(t, "42", id.String())
		assert.Equal(t, uint32(42), id.UInt32())
	})
}
