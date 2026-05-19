package character_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestService spins up an in-memory database, pre-registers the labels
// used by the tests and returns a ready-to-use character service.
func newTestService(t *testing.T) (*character.Service, *database.DB, map[string]uint32) {
	t.Helper()

	db := database.NewDB()
	svc := character.NewService(db)

	labels := map[string]uint32{}
	for _, l := range []string{
		"hero", "villain", "wizard",
		"noun-man", "noun-woman", "noun-mage",
		"adj-brave", "adj-evil", "adj-wise",
	} {
		id, err := db.CreateLabel(l)
		require.NoError(t, err, "creating label %q", l)
		labels[l] = id
	}

	return svc, db, labels
}

// mustCreate builds a character via the build func and persists it.
func mustCreate(t *testing.T, db *database.DB, build func(c *character.Character)) *character.Character {
	t.Helper()

	c := character.New()
	build(c)

	id, err := db.Create(c)
	require.NoError(t, err)
	require.NotZero(t, id)

	return c
}

func TestNewAndAccessors(t *testing.T) {
	c := character.New()

	assert.NotNil(t, c)
	assert.Equal(t, kind.Character, c.GetKind())
	assert.Zero(t, c.GetID())
	assert.Zero(t, c.GetLabelID())
	assert.False(t, c.Human)
	assert.False(t, c.Created)

	c.SetID(42)
	assert.Equal(t, uint32(42), c.GetID())

	c.SetLabelID(7)
	assert.Equal(t, uint32(7), c.GetLabelID())
}

func TestQuery(t *testing.T) {
	svc, db, labels := newTestService(t)

	hero := mustCreate(t, db, func(c *character.Character) {
		c.LabelID = labels["hero"]
		c.NounID = labels["noun-man"]
		c.AdjectiveID = labels["adj-brave"]
		c.Description = "the hero"
		c.LocationID = 1
		c.Human = true
	})
	villain := mustCreate(t, db, func(c *character.Character) {
		c.LabelID = labels["villain"]
		c.NounID = labels["noun-man"]
		c.AdjectiveID = labels["adj-evil"]
		c.Description = "the villain"
		c.LocationID = 2
	})
	wizard := mustCreate(t, db, func(c *character.Character) {
		c.LabelID = labels["wizard"]
		c.NounID = labels["noun-mage"]
		c.AdjectiveID = labels["adj-wise"]
		c.Description = "the wizard"
		c.LocationID = 3
	})

	t.Run("Count", func(t *testing.T) {
		assert.Equal(t, 3, svc.Get().Count())
	})

	t.Run("WithHuman true returns only the human", func(t *testing.T) {
		assert.Equal(t, 1, svc.Get().WithHuman(true).Count())

		got, err := svc.Get().WithHuman(true).First()
		require.NoError(t, err)
		assert.Equal(t, hero.GetID(), got.GetID())
	})

	t.Run("WithHuman false returns only the non-humans", func(t *testing.T) {
		assert.Equal(t, 2, svc.Get().WithHuman(false).Count())
	})

	t.Run("WithID", func(t *testing.T) {
		got, err := svc.Get().WithID(villain.GetID()).First()
		require.NoError(t, err)
		assert.Equal(t, "the villain", got.Description)
	})

	t.Run("WithLabel", func(t *testing.T) {
		got, err := svc.Get().WithLabel("wizard").First()
		require.NoError(t, err)
		assert.Equal(t, wizard.GetID(), got.GetID())
	})

	t.Run("WithLabelID", func(t *testing.T) {
		got, err := svc.Get().WithLabelID(labels["hero"]).First()
		require.NoError(t, err)
		assert.Equal(t, hero.GetID(), got.GetID())
	})

	t.Run("WithNoID excludes", func(t *testing.T) {
		// Two non-hero characters left when excluding the hero id.
		assert.Equal(t, 2, svc.Get().WithNoID(hero.GetID()).Count())
	})

	t.Run("Exists by label", func(t *testing.T) {
		assert.True(t, svc.Get().WithLabel("hero").Exists())
		assert.False(t, svc.Get().WithLabel("nonexistent").Exists())
	})

	t.Run("First missing returns error", func(t *testing.T) {
		_, err := svc.Get().WithID(99999).First()
		require.Error(t, err)
	})
}

func TestServices(t *testing.T) {
	t.Run("Count is zero on an empty service", func(t *testing.T) {
		svc, _, _ := newTestService(t)
		assert.Zero(t, svc.Count())
	})

	t.Run("HasHuman and GetHuman", func(t *testing.T) {
		svc, db, labels := newTestService(t)

		// No human yet
		assert.False(t, svc.HasHuman())

		_, err := svc.GetHuman()
		require.Error(t, err)

		// Add a non-human first to make sure WithHuman(true) discriminates
		mustCreate(t, db, func(c *character.Character) {
			c.LabelID = labels["villain"]
			c.NounID = labels["noun-man"]
			c.AdjectiveID = labels["adj-evil"]
			c.Description = "the villain"
			c.LocationID = 2
		})

		assert.False(t, svc.HasHuman())

		// Now add the human
		hero := mustCreate(t, db, func(c *character.Character) {
			c.LabelID = labels["hero"]
			c.NounID = labels["noun-man"]
			c.AdjectiveID = labels["adj-brave"]
			c.Description = "the hero"
			c.LocationID = 1
			c.Human = true
		})

		assert.True(t, svc.HasHuman())

		got, err := svc.GetHuman()
		require.NoError(t, err)
		assert.Equal(t, hero.GetID(), got.GetID())
	})

	t.Run("GetByLabel", func(t *testing.T) {
		svc, db, labels := newTestService(t)

		wizard := mustCreate(t, db, func(c *character.Character) {
			c.LabelID = labels["wizard"]
			c.NounID = labels["noun-mage"]
			c.AdjectiveID = labels["adj-wise"]
			c.Description = "the wizard"
			c.LocationID = 3
		})

		got, err := svc.GetByLabel("wizard")
		require.NoError(t, err)
		assert.Equal(t, wizard.GetID(), got.GetID())

		_, err = svc.GetByLabel("nobody")
		require.Error(t, err)
	})
}

func TestValidateAll(t *testing.T) {
	validCharBuilder := func(db *database.DB, labels map[string]uint32) func(label, noun, adj string, human bool) {
		return func(label, noun, adj string, human bool) {
			t.Helper()
			mustCreate(t, db, func(c *character.Character) {
				c.LabelID = labels[label]
				c.NounID = labels[noun]
				c.AdjectiveID = labels[adj]
				c.Description = "ok"
				c.LocationID = 1
				c.Human = human
			})
		}
	}

	t.Run("no human reports ErrNoHuman", func(t *testing.T) {
		svc, db, labels := newTestService(t)
		create := validCharBuilder(db, labels)

		create("villain", "noun-man", "adj-evil", false)

		errs := svc.ValidateAll()
		require.Len(t, errs, 1)
		require.ErrorIs(t, errs[0], character.ErrNoHuman)
	})

	t.Run("empty database reports ErrNoHuman", func(t *testing.T) {
		svc, _, _ := newTestService(t)

		errs := svc.ValidateAll()
		require.Len(t, errs, 1)
		require.ErrorIs(t, errs[0], character.ErrNoHuman)
	})

	t.Run("exactly one human is valid", func(t *testing.T) {
		svc, db, labels := newTestService(t)
		create := validCharBuilder(db, labels)

		create("hero", "noun-man", "adj-brave", true)
		create("villain", "noun-man", "adj-evil", false)

		assert.Empty(t, svc.ValidateAll())
	})

	t.Run("two humans reports ErrOnlyOneHuman", func(t *testing.T) {
		svc, db, labels := newTestService(t)
		create := validCharBuilder(db, labels)

		create("hero", "noun-man", "adj-brave", true)
		create("wizard", "noun-mage", "adj-wise", true)

		errs := svc.ValidateAll()

		var hit bool
		for _, e := range errs {
			if errors.Is(e, character.ErrOnlyOneHuman) {
				hit = true

				break
			}
		}
		require.True(t, hit, "expected ErrOnlyOneHuman in %v", errs)
	})

	t.Run("missing required fields are reported", func(t *testing.T) {
		svc, db, labels := newTestService(t)

		// Valid human so the "no human" rule does not pollute the result.
		mustCreate(t, db, func(c *character.Character) {
			c.LabelID = labels["hero"]
			c.NounID = labels["noun-man"]
			c.AdjectiveID = labels["adj-brave"]
			c.Description = "the hero"
			c.LocationID = 1
			c.Human = true
		})

		// Incomplete character: missing Description and LocationID.
		mustCreate(t, db, func(c *character.Character) {
			c.LabelID = labels["villain"]
			c.NounID = labels["noun-man"]
			c.AdjectiveID = labels["adj-evil"]
		})

		errs := svc.ValidateAll()
		require.NotEmpty(t, errs)

		var sawRequired bool
		for _, e := range errs {
			msg := e.Error()
			if strings.Contains(msg, "Description is required") ||
				strings.Contains(msg, "LocationID is required") {
				sawRequired = true

				break
			}
		}
		assert.True(t, sawRequired, "expected a required-field error, got %v", errs)
	})
}
