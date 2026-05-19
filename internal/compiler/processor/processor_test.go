package processor_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/processor"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setup builds a fresh adventure + status pair with the given section preset.
func setup(t *testing.T, section kind.Kind) (*adventure.Adventure, *status.Status) {
	t.Helper()

	a := adventure.New()
	st := status.New(a.DB)
	st.Section = section

	return a, st
}

// mk builds a Line with the given text.
func mk(text string) line.Line {
	return line.New(text, 1)
}

// seedWord stores a word in the adventure so downstream processors can resolve it.
func seedWord(t *testing.T, a *adventure.Adventure, wt word.WordType, syns ...string) *word.Word {
	t.Helper()

	labelID, err := a.DB.CreateLabel(syns[0])
	require.NoError(t, err)

	w := word.New(labelID, wt, syns...)
	_, err = a.Words.Create(w)
	require.NoError(t, err)

	return w
}

// seedLocation creates a location bypassing the processor (for setting up
// references that other processors will resolve via labels).
func seedLocation(t *testing.T, a *adventure.Adventure, label string) *location.Location {
	t.Helper()

	labelID, err := a.DB.CreateLabel(label)
	require.NoError(t, err)

	loc := location.New()
	loc.LabelID = labelID
	loc.Title = label
	loc.Description = label

	_, err = a.Locations.Create(loc)
	require.NoError(t, err)

	return loc
}

func TestProcessLineDispatch(t *testing.T) {
	t.Run("None section is a no-op", func(t *testing.T) {
		a, st := setup(t, kind.None)
		require.NoError(t, processor.ProcessLine(mk("anything"), st, a))
	})

	t.Run("routes to Param", func(t *testing.T) {
		a, st := setup(t, kind.Param)
		require.NoError(t, processor.ProcessLine(mk(`title: "Some title"`), st, a))
		assert.True(t, a.Config.Get().WithLabel("title").Exists())
	})

	t.Run("routes to Word", func(t *testing.T) {
		a, st := setup(t, kind.Word)
		require.NoError(t, processor.ProcessLine(mk("verb: north, n"), st, a))
		assert.True(t, a.Words.Get().WithLabel("north").Exists())
	})
}

func TestReadConfig(t *testing.T) {
	t.Run("valid config sets the param", func(t *testing.T) {
		a, st := setup(t, kind.Param)
		require.NoError(t, processor.ProcessLine(mk(`author: "queru"`), st, a))

		p, err := a.Config.Get().WithLabel("author").First()
		require.NoError(t, err)
		assert.Equal(t, "queru", p.V)
	})

	t.Run("garbage line returns error", func(t *testing.T) {
		a, st := setup(t, kind.Param)
		require.Error(t, processor.ProcessLine(mk("just garbage"), st, a))
	})

	t.Run("unknown label returns error", func(t *testing.T) {
		a, st := setup(t, kind.Param)
		require.Error(t, processor.ProcessLine(mk(`nonExistentField: "x"`), st, a))
	})
}

func TestReadVar(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		a, st := setup(t, kind.Variable)
		require.NoError(t, processor.ProcessLine(mk("count = 42"), st, a))

		v, err := a.Variables.Get().WithLabel("count").First()
		require.NoError(t, err)
		assert.Equal(t, 42, v.Int())
		assert.True(t, v.IsInt())
	})

	t.Run("float", func(t *testing.T) {
		a, st := setup(t, kind.Variable)
		require.NoError(t, processor.ProcessLine(mk("ratio = 1.5"), st, a))

		v, err := a.Variables.Get().WithLabel("ratio").First()
		require.NoError(t, err)
		assert.InDelta(t, 1.5, v.Float(), 1e-9)
		assert.True(t, v.IsFloat())
	})

	t.Run("bool true", func(t *testing.T) {
		a, st := setup(t, kind.Variable)
		require.NoError(t, processor.ProcessLine(mk("isOn = true"), st, a))

		v, err := a.Variables.Get().WithLabel("isOn").First()
		require.NoError(t, err)
		assert.True(t, v.Bool())
	})

	t.Run("bool false", func(t *testing.T) {
		a, st := setup(t, kind.Variable)
		require.NoError(t, processor.ProcessLine(mk("isOff = false"), st, a))

		v, err := a.Variables.Get().WithLabel("isOff").First()
		require.NoError(t, err)
		assert.False(t, v.Bool())
	})

	t.Run("string", func(t *testing.T) {
		a, st := setup(t, kind.Variable)
		require.NoError(t, processor.ProcessLine(mk(`name = "Alice"`), st, a))

		v, err := a.Variables.Get().WithLabel("name").First()
		require.NoError(t, err)
		assert.Equal(t, "Alice", v.String())
	})

	t.Run("garbage returns error", func(t *testing.T) {
		a, st := setup(t, kind.Variable)
		require.Error(t, processor.ProcessLine(mk("not a var"), st, a))
	})
}

func TestReadWord(t *testing.T) {
	t.Run("creates a verb with synonyms", func(t *testing.T) {
		a, st := setup(t, kind.Word)
		require.NoError(t, processor.ProcessLine(mk("verb: north, n"), st, a))

		w, err := a.Words.Get().WithLabel("north").First()
		require.NoError(t, err)
		assert.Equal(t, word.Verb, w.Type)
		assert.ElementsMatch(t, []string{"north", "n"}, w.Synonyms)
	})

	t.Run("creates a noun", func(t *testing.T) {
		a, st := setup(t, kind.Word)
		require.NoError(t, processor.ProcessLine(mk("noun: sword, blade"), st, a))

		w, err := a.Words.Get().WithLabel("sword").First()
		require.NoError(t, err)
		assert.Equal(t, word.Noun, w.Type)
	})

	t.Run("garbage returns error", func(t *testing.T) {
		a, st := setup(t, kind.Word)
		require.Error(t, processor.ProcessLine(mk("just words without colon"), st, a))
	})

	t.Run("unknown word type returns error", func(t *testing.T) {
		a, st := setup(t, kind.Word)
		require.Error(t, processor.ProcessLine(mk("xyz: foo, bar"), st, a))
	})
}

func TestReadMessage(t *testing.T) {
	t.Run("creates a single-line message", func(t *testing.T) {
		a, st := setup(t, kind.Message)
		require.NoError(t, processor.ProcessLine(mk(`greeting: "Hello there"`), st, a))

		// Messages buffer current entity until a new one or a section change.
		// Flush it explicitly.
		require.True(t, st.SaveCurrentStoreable().IsOK())

		m, err := a.Messages.Get().WithLabel("greeting").First()
		require.NoError(t, err)
		assert.NotZero(t, m.GetID())
		assert.Equal(t, "Hello there", m.Text)
	})

	t.Run("plural variants of the same label merge into one message", func(t *testing.T) {
		a, st := setup(t, kind.Message)

		require.NoError(t, processor.ProcessLine(mk(`order_count.zero: "No orders yet."`), st, a))
		require.NoError(t, processor.ProcessLine(mk(`order_count.one: "One order."`), st, a))
		require.NoError(t, processor.ProcessLine(mk(`order_count.many: "_ orders."`), st, a))

		require.True(t, st.SaveCurrentStoreable().IsOK())

		// Only one message should exist for "order_count".
		assert.Equal(t, 1, a.Messages.Get().WithLabel("order_count").Count())

		m, err := a.Messages.Get().WithLabel("order_count").First()
		require.NoError(t, err)
		assert.Equal(t, "No orders yet.", m.Text, "zero variant lives in Text")
		assert.True(t, m.IsPluralized(), "one/many variants must be present")
	})

	t.Run("different labels stay as separate messages", func(t *testing.T) {
		a, st := setup(t, kind.Message)

		require.NoError(t, processor.ProcessLine(mk(`hello: "Hi"`), st, a))
		require.NoError(t, processor.ProcessLine(mk(`bye: "Bye"`), st, a))

		require.True(t, st.SaveCurrentStoreable().IsOK())

		assert.True(t, a.Messages.Get().WithLabel("hello").Exists())
		assert.True(t, a.Messages.Get().WithLabel("bye").Exists())
	})

	t.Run("garbage returns error", func(t *testing.T) {
		a, st := setup(t, kind.Message)
		require.Error(t, processor.ProcessLine(mk("not-a-message"), st, a))
	})
}

func TestReadLocation(t *testing.T) {
	t.Run("full lifecycle of a location", func(t *testing.T) {
		a, st := setup(t, kind.Location)
		seedWord(t, a, word.Verb, "north", "n")

		require.NoError(t, processor.ProcessLine(mk("entrance:"), st, a))
		require.NoError(t, processor.ProcessLine(mk(`	title: "The entrance"`), st, a))
		require.NoError(t, processor.ProcessLine(mk(`	desc: "A long hallway"`), st, a))
		require.NoError(t, processor.ProcessLine(mk("	exits: north hall"), st, a))

		require.True(t, st.SaveCurrentStoreable().IsOK())

		got, err := a.Locations.Get().WithLabel("entrance").First()
		require.NoError(t, err)
		assert.Equal(t, "The entrance", got.Title)
		assert.Equal(t, "A long hallway", got.Description)

		// The seeded verb should now be flagged as a connection word.
		w, err := a.Words.Get().WithLabel("north").First()
		require.NoError(t, err)
		assert.True(t, w.IsConnection)
	})

	t.Run("exit referencing an unknown word fails", func(t *testing.T) {
		a, st := setup(t, kind.Location)

		require.NoError(t, processor.ProcessLine(mk("entrance:"), st, a))
		require.Error(t, processor.ProcessLine(mk("	exits: nowhere somewhere"), st, a))
	})

	t.Run("garbage returns error", func(t *testing.T) {
		a, st := setup(t, kind.Location)
		require.Error(t, processor.ProcessLine(mk("nonsense line"), st, a))
	})
}

func TestReadItem(t *testing.T) {
	t.Run("full lifecycle of an item", func(t *testing.T) {
		a, st := setup(t, kind.Item)
		seedWord(t, a, word.Noun, "sword")
		seedWord(t, a, word.Adjective, "shiny")

		require.NoError(t, processor.ProcessLine(mk("excalibur: sword shiny"), st, a))
		require.NoError(t, processor.ProcessLine(mk(`	desc: "A shiny sword"`), st, a))
		require.NoError(t, processor.ProcessLine(mk("	is wearable"), st, a))
		require.NoError(t, processor.ProcessLine(mk("	is container"), st, a))
		require.NoError(t, processor.ProcessLine(mk("	is created"), st, a))
		require.NoError(t, processor.ProcessLine(mk("	is at hall"), st, a))
		require.NoError(t, processor.ProcessLine(mk("	weight 5"), st, a))
		require.NoError(t, processor.ProcessLine(mk("	max weight 50"), st, a))

		require.True(t, st.SaveCurrentStoreable().IsOK())

		got, err := a.Items.Get().WithLabel("excalibur").First()
		require.NoError(t, err)
		assert.Equal(t, "A shiny sword", got.Description)
		assert.True(t, got.Wearable)
		assert.True(t, got.Container)
		assert.True(t, got.Created)
		assert.Equal(t, 5, got.Weight)
		assert.Equal(t, 50, got.MaxWeight)
		assert.NotZero(t, got.At, "is at should resolve to a label id")
	})

	t.Run("unknown noun returns error", func(t *testing.T) {
		a, st := setup(t, kind.Item)
		seedWord(t, a, word.Adjective, "shiny")

		require.Error(t, processor.ProcessLine(mk("excalibur: unknown shiny"), st, a))
	})

	t.Run("unknown adjective returns error", func(t *testing.T) {
		a, st := setup(t, kind.Item)
		seedWord(t, a, word.Noun, "sword")

		require.Error(t, processor.ProcessLine(mk("excalibur: sword unknown"), st, a))
	})
}

func TestReadCharacter(t *testing.T) {
	t.Run("full lifecycle of a character", func(t *testing.T) {
		a, st := setup(t, kind.Character)
		seedWord(t, a, word.Noun, "hero")
		seedWord(t, a, word.Adjective, "brave")
		hall := seedLocation(t, a, "hall")

		require.NoError(t, processor.ProcessLine(mk("john: hero brave"), st, a))
		require.NoError(t, processor.ProcessLine(mk(`	desc: "The hero"`), st, a))
		require.NoError(t, processor.ProcessLine(mk("	is human"), st, a))
		require.NoError(t, processor.ProcessLine(mk("	is created"), st, a))
		require.NoError(t, processor.ProcessLine(mk("	is at hall"), st, a))

		require.True(t, st.SaveCurrentStoreable().IsOK())

		got, err := a.Characters.Get().WithLabel("john").First()
		require.NoError(t, err)
		assert.Equal(t, "The hero", got.Description)
		assert.True(t, got.Human)
		assert.True(t, got.Created)
		assert.Equal(t, hall.ID, got.LocationID)
	})

	t.Run("second human declaration fails", func(t *testing.T) {
		a, st := setup(t, kind.Character)
		seedWord(t, a, word.Noun, "hero")
		seedWord(t, a, word.Noun, "wizard")
		seedWord(t, a, word.Adjective, "brave")
		seedWord(t, a, word.Adjective, "wise")

		// First character, marked as human.
		require.NoError(t, processor.ProcessLine(mk("john: hero brave"), st, a))
		require.NoError(t, processor.ProcessLine(mk("	is human"), st, a))

		// Header of the second character flushes the first.
		require.NoError(t, processor.ProcessLine(mk("merlin: wizard wise"), st, a))

		// Attempting a second human now must fail.
		require.Error(t, processor.ProcessLine(mk("	is human"), st, a))
	})

	t.Run("is at unknown location fails", func(t *testing.T) {
		a, st := setup(t, kind.Character)
		seedWord(t, a, word.Noun, "hero")
		seedWord(t, a, word.Adjective, "brave")

		require.NoError(t, processor.ProcessLine(mk("john: hero brave"), st, a))
		require.Error(t, processor.ProcessLine(mk("	is at nowhere"), st, a))
	})

	t.Run("unknown noun returns error", func(t *testing.T) {
		a, st := setup(t, kind.Character)
		seedWord(t, a, word.Adjective, "brave")

		require.Error(t, processor.ProcessLine(mk("john: unknown brave"), st, a))
	})
}

func TestReadBlob(t *testing.T) {
	t.Run("happy path loads a file into a blob", func(t *testing.T) {
		dir := t.TempDir()
		assetDir := filepath.Join(dir, "gfx")
		require.NoError(t, os.MkdirAll(assetDir, 0o755))

		assetPath := filepath.Join(assetDir, "logo.png")
		require.NoError(t, os.WriteFile(assetPath, []byte{0x89, 'P', 'N', 'G'}, 0o644))

		a, st := setup(t, kind.Blob)
		// readBlob resolves the path against status.CurrentFilename's directory,
		// and the Blob regex requires at least one slash in the filename.
		st.PushFilename(filepath.Join(dir, "adventure.adv"))

		require.NoError(t, processor.ProcessLine(mk("logo: gfx/logo.png"), st, a))

		got, err := a.Blobs.Get().WithLabel("logo").First()
		require.NoError(t, err)
		assert.Equal(t, []byte{0x89, 'P', 'N', 'G'}, got.Data)
		assert.NotEmpty(t, got.Mime)
	})

	t.Run("garbage returns error", func(t *testing.T) {
		a, st := setup(t, kind.Blob)
		require.Error(t, processor.ProcessLine(mk("not-a-blob-decl"), st, a))
	})

	t.Run("missing file returns error", func(t *testing.T) {
		a, st := setup(t, kind.Blob)
		st.PushFilename(filepath.Join(t.TempDir(), "adventure.adv"))

		require.Error(t, processor.ProcessLine(mk("logo: gfx/missing.png"), st, a))
	})
}
