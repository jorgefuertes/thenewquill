package compiler_test

import (
	"bytes"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/compiler"
	"github.com/jorgefuertes/thenewquill/internal/compiler/db"

	"github.com/stretchr/testify/require"
)

func TestCompilerDB(t *testing.T) {
	const srcFilename = "src/happy/test.adv"

	var a *adventure.Adventure
	var database *db.DB
	var save *bytes.Buffer
	var load *bytes.Buffer

	var a2 *adventure.Adventure
	var database2 *db.DB

	t.Run("compile", func(t *testing.T) {
		var err error

		a, err = compiler.Compile(srcFilename)
		require.NoError(t, err)
		require.NotNil(t, a)

		t.Run("adventure to DB", func(t *testing.T) {
			database = db.New()
			a.Export(database)
			require.NotNil(t, database)
			require.NotEmpty(t, database.Headers)
			require.NotEmpty(t, database.Records)
			for i, r := range database.Records {
				require.NotEmpty(t, r.Section, "section empty in reg %d", i)
				require.NotEmpty(t, r.Label, "label empty in reg %d", i)
				require.NotEmpty(t, r.Fields, "fields empty in reg %d", i)
			}

			t.Run("DB to file", func(t *testing.T) {
				save = bytes.NewBuffer(nil)
				err := database.Save(save)
				require.NoError(t, err)
				require.NotNil(t, save)
				require.NotZero(t, save.Len())
				load = bytes.NewBuffer(save.Bytes())

				t.Run("file to DB", func(t *testing.T) {
					database2 = db.New()
					err := database2.Load(load)
					require.NoError(t, err)

					h1, err := database.Hash()
					require.NoError(t, err)
					h2, err := database2.Hash()
					require.NoError(t, err)
					require.Equal(t, h1, h2)

					require.Equal(t, database.Headers, database2.Headers)
					for i, r := range database.Records {
						require.Equal(t, r, database2.Records[i])
					}

					t.Run("adventure from DB", func(t *testing.T) {
						a2 = adventure.New()
						err := a2.Import(database2)
						require.NoError(t, err)

						assertEqualAventures(t, a, a2)
					})
				})
			})
		})
	})
}
