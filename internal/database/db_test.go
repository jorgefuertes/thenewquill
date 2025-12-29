package database_test

import (
	"fmt"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestItem struct {
	ID      uint32
	LabelID uint32
	Title   string
	At      uint32
	OK      bool
	NOOK    bool
	Weight  int
	Names   []string
	Numbers []int
}

const numOfItems = 1000

func TestDatabase(t *testing.T) {
	db := database.NewDB()
	require.NotEmpty(t, db)

	t.Run("write", func(t *testing.T) {
		fillDatabase(t, db)
		assert.Equal(t, numOfItems, db.CountRecords())

		t.Run("reset", func(t *testing.T) {
			db.ResetDB()
			assert.Equal(t, 0, db.CountRecords())
		})
	})

	t.Run("retrieve", func(t *testing.T) {
		t.Cleanup(db.ResetDB)
		fillDatabase(t, db)

		t.Run("get by ID", func(t *testing.T) {
			for i := range numOfItems {
				id := uint32(i + 1)
				var item TestItem

				err := db.Get(id, &item)
				require.NoError(t, err)
				require.Equal(t, id, item.ID)
				require.NotZero(t, item.LabelID)
				require.Equal(t, titleForIteration(id-1), item.Title)
				require.EqualValues(t, id-1, item.At)
				require.True(t, item.OK)
				require.False(t, item.NOOK)
				require.EqualValues(t, id+100-1, item.Weight)
				require.Equal(t, []string{"one", "two", "three"}, item.Names)
				require.Equal(t, []int{1, 2, 3}, item.Numbers)
			}
		})

		t.Run("Get by label", func(t *testing.T) {
			i := rand.Uint32N(numOfItems)
			id := uint32(i + 3)

			var item TestItem
			err := db.GetByLabel(labelForIteration(i), &item)
			require.NoError(t, err)
			require.EqualValues(t, id, item.LabelID, "label ID %d is not equal to %d", i, item.LabelID)
		})
	})

	t.Run("IO", func(t *testing.T) {
		t.Cleanup(func() {
			assert.NoError(t, os.Remove("/tmp/test.db"))
			db.ResetDB()
		})

		fillDatabase(t, db)

		t.Run("export", func(t *testing.T) {
			filename := "/tmp/test.db"
			count, size, err := db.Export(filename)
			require.NoError(t, err)
			require.NotZero(t, count)
			t.Logf("%d bytes writen, %d bytes final size", count, size)

			t.Run("import", func(t *testing.T) {
				err := db.Import(filename)
				require.NoError(t, err)
			})
		})
	})
}

func fillDatabase(t *testing.T, db *database.DB) {
	t.Helper()

	for i := range numOfItems {
		labelID, err := db.CreateLabel(labelForIteration(uint32(i)))
		require.NoError(t, err)
		require.NotZero(t, labelID)
		require.EqualValues(t, i+3, labelID)

		item := TestItem{
			ID:      0,
			LabelID: labelID,
			Title:   titleForIteration(uint32(i)),
			At:      uint32(i),
			OK:      true,
			NOOK:    false,
			Weight:  i + 100,
			Names:   []string{"one", "two", "three"},
			Numbers: []int{1, 2, 3},
		}

		id, err := db.Create(&item)
		require.NoError(t, err)
		require.NotZero(t, id)
		require.EqualValues(t, i+1, id)
	}
}

func titleForIteration(i uint32) string {
	return fmt.Sprintf("Title for iteration %d", i)
}

func labelForIteration(i uint32) string {
	return fmt.Sprintf("label-it-%d", i)
}
