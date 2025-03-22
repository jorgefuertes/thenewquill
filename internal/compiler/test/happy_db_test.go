package compiler_test

// func TestCompilerDB(t *testing.T) {
// 	const srcFilename = "src/happy/test.adv"

// 	var a *adventure.Adventure
// 	var database *db.DB
// 	var save *bytes.Buffer
// 	var load *bytes.Buffer

// 	var a2 *adventure.Adventure
// 	var database2 *db.DB

// 	t.Run("compile", func(t *testing.T) {
// 		var err error

// 		a, err = compiler.Compile(srcFilename)
// 		require.NoError(t, err)
// 		require.NotNil(t, a)

// 		t.Run("adventure to DB", func(t *testing.T) {
// 			database = db.NewDB()
// 			a.Export(database)
// 			require.NotNil(t, database)
// 			require.NotEmpty(t, database.GetHeaders())
// 			require.NotEmpty(t, database.GetRegs())
// 			for i, r := range database.GetRegs() {
// 				require.NotEmpty(t, r.Section, "section empty in reg %d", i)
// 				require.NotEmpty(t, r.Label, "label empty in reg %d", i)
// 				require.NotEmpty(t, r.Fields, "fields empty in reg %d", i)
// 			}

// 			t.Run("DB to file", func(t *testing.T) {
// 				save = bytes.NewBuffer(nil)
// 				err := database.Write(save)
// 				require.NoError(t, err)
// 				require.NotNil(t, save)
// 				require.NotZero(t, save.Len())
// 				load = bytes.NewBuffer(save.Bytes())

// 				t.Run("file to DB", func(t *testing.T) {
// 					database2 = db.NewDB()
// 					err := database2.Load(load)
// 					require.NoError(t, err)

// 					require.Equal(t, database.Hash(), database2.Hash())
// 					require.Equal(t, database.GetHeaders(), database2.GetHeaders())
// 					for i, r := range database.GetRegs() {
// 						require.Equal(t, r, database2.GetRegs()[i])
// 					}

// 					t.Run("adventure from DB", func(t *testing.T) {
// 						a2 = adventure.New()
// 						err := a2.Import(database2)
// 						require.NoError(t, err)

// 						require.Equal(t, a.Config, a2.Config)
// 						require.Equal(t, a.Vars, a2.Vars)
// 						require.Equal(t, a.Words, a2.Words)
// 						require.Equal(t, a.Messages, a2.Messages)
// 						require.Equal(t, a.Locations, a2.Locations)
// 						require.Equal(t, a.Items, a2.Items)
// 						require.Equal(t, a.Chars, a2.Chars)
// 					})
// 				})
// 			})
// 		})
// 	})
// }
