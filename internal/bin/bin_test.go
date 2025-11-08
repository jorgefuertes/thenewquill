package bin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/variable"
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

	currentID := db.ID(5)

	t.Run("write character", func(t *testing.T) {
		t.Run("write label", func(t *testing.T) {
			l := db.Label{ID: currentID, Name: "human"}
			err := b.WriteStoreable(l)
			require.NoError(t, err)
		})

		c := character.Character{
			ID:          currentID,
			NounID:      2,
			AdjectiveID: 3,
			Description: "description",
			LocationID:  4,
			Created:     false,
			Human:       true,
		}

		currentID++

		err := b.WriteStoreable(c)
		require.NoError(t, err)
	})

	testVars := []struct {
		name string
		val  any
	}{
		{"true", true},
		{"false", false},
		{"number", 10},
		{"number2", 20},
		{"aFloat", 1.5},
		{"name", "The New Quill Adventure Writing System"},
		{"hello", "Hello, _.\nWelcome to _.\n"},
		{"gradas.people", 500},
		{"subasta.running", true},
		{"via.open", true},
		{"ant-on", false},
		{"player.health", 100},
		{"enano.patience", 255},
		{"enano.death", false},
		{"elfo.hidden", true},
	}

	t.Run("write variables", func(t *testing.T) {
		for _, tc := range testVars {
			t.Run(tc.name, func(t *testing.T) {
				l := db.Label{ID: currentID, Name: tc.name}
				err := b.WriteStoreable(l)
				require.NoError(t, err)

				v := variable.New(l.ID, tc.val)

				currentID++

				err = b.WriteStoreable(v)
				require.NoError(t, err)
			})
		}
	})

	t.Run("save to disk", func(t *testing.T) {
		n, err := b.Save(filename)
		require.NoError(t, err)
		require.NotZero(t, n)
	})

	t.Run("load from disk", func(t *testing.T) {
		err := b.Load(filename)
		require.NoError(t, err)
	})

	currentID = db.ID(5)

	t.Run("read character", func(t *testing.T) {
		t.Run("read label", func(t *testing.T) {
			l, err := b.ReadStoreable()
			require.NoError(t, err)
			require.IsType(t, db.Label{}, l)
			require.Equal(t, currentID, l.GetID())
		})

		s, err := b.ReadStoreable()
		require.NoError(t, err)
		require.IsType(t, character.Character{}, s)
		require.Equal(t, currentID, s.GetID())

		c, ok := s.(character.Character)
		require.True(t, ok)
		require.Equal(t, db.ID(2), c.NounID)
		require.Equal(t, db.ID(3), c.AdjectiveID)
		require.Equal(t, "description", c.Description)
		require.Equal(t, db.ID(4), c.LocationID)
		require.Equal(t, false, c.Created)
		require.Equal(t, true, c.Human)

		currentID++
	})

	t.Run("read variables", func(t *testing.T) {
		for _, tc := range testVars {
			t.Run(tc.name, func(t *testing.T) {
				l, err := b.ReadStoreable()
				require.NoError(t, err)
				require.Equal(t, kind.Label, kind.KindOf(l))
				require.Equal(t, currentID, l.GetID())
				label, ok := l.(db.Label)
				require.True(t, ok, "label should be a db.Label and not a %T", l)
				require.Equal(t, tc.name, label.Name)

				st, err := b.ReadStoreable()
				require.NoError(t, err)
				require.Equal(t, kind.Variable, kind.KindOf(st))
				require.Equal(t, currentID, st.GetID())

				v2, ok := st.(variable.Variable)
				require.True(t, ok)

				switch val := tc.val.(type) {
				case bool:
					require.Equal(t, val, v2.Bool())
				case int:
					require.Equal(t, val, v2.Int())
				case float32, float64:
					require.Equal(t, val, v2.Float())
				case string:
					require.Equal(t, val, v2.String())
				default:
					t.Errorf("unexpected type: %T", tc.val)
				}

				currentID++
			})
		}
	})
}
