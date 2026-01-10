package database

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportImport(t *testing.T) {
	randomUint32 := func() uint32 {
		return uint32(rand.Intn(math.MaxInt32))
	}

	randomInt := func() int {
		return rand.Int()
	}

	randomBool := func() bool {
		return rand.Intn(2) == 1
	}

	randomStrings := func(max int) []string {
		count := rand.Intn(max) + 1
		strs := make([]string, count)
		for i := 0; i < count; i++ {
			strs[i] = fmt.Sprintf("str_%d", rand.Intn(100000))
		}

		return strs
	}

	randomInts := func(max int) []int {
		count := rand.Intn(max) + 1
		ints := make([]int, count)
		for i := 0; i < count; i++ {
			ints[i] = rand.Intn(100000)
		}

		return ints
	}

	createRandomRecords := func(db *DB, count int) {
		for i := 0; i < count; i++ {
			labelID, err := db.CreateLabel(fmt.Sprintf("test_%d", i))
			require.NoError(t, err)

			ti := &testItem{
				LabelID: labelID,
				Title:   fmt.Sprintf("Test %d", i),
				At:      randomUint32(),
				OK:      randomBool(),
				NOOK:    randomBool(),
				Weight:  randomInt(),
				Names:   randomStrings(5),
				Numbers: randomInts(5),
			}

			id, err := db.Create(ti)
			require.NoError(t, err)
			require.NotZero(t, id)
		}
	}

	tmp := t.TempDir()
	origDB := NewDB()
	dstDB := NewDB()

	t.Run("DB", func(t *testing.T) {
		createRandomRecords(origDB, 100)

		t.Run("Export", func(t *testing.T) {
			bSent, bFile, err := origDB.Export(tmp + "/test.db")
			require.NoError(t, err)
			assert.NotZero(t, bSent)
			assert.NotZero(t, bFile)

			t.Run("Import", func(t *testing.T) {
				err := dstDB.Import(tmp + "/test.db")
				require.NoError(t, err)

				t.Run("Compare", func(t *testing.T) {
					assert.EqualValues(t, origDB.CountLabels(), dstDB.CountLabels())
					assert.EqualValues(t, origDB.CountRecords(), dstDB.CountRecords())

					for id, origLabel := range origDB.labels {
						dstLabel, ok := dstDB.labels[id]
						require.True(t, ok, "label %d not found in dstDB", id)

						assert.EqualValues(t, origLabel, dstLabel)
					}

					for id, origRec := range origDB.data {
						dstRec, ok := dstDB.data[id]
						require.True(t, ok, "record %d not found in dstDB", id)

						assert.EqualValues(t, origRec.LabelID, dstRec.LabelID)
						assert.EqualValues(t, origRec.Kind, dstRec.Kind)
						assert.EqualValues(t, origRec.Data, dstRec.Data)
					}
				})
			})
		})
	})
}
