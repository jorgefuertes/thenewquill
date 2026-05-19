package status_test

import (
	"path/filepath"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/item"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/compiler/status"
	"github.com/jorgefuertes/thenewquill/internal/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newStatus(t *testing.T) (*status.Status, *database.DB) {
	t.Helper()

	db := database.NewDB()
	return status.New(db), db
}

func TestNew(t *testing.T) {
	s, _ := newStatus(t)

	assert.Equal(t, kind.None, s.Section)
	assert.Empty(t, s.Stack)
	assert.False(t, s.Comment.IsOn())
	assert.False(t, s.MultiLine.IsOn())
	assert.False(t, s.HasCurrent())
	assert.Empty(t, s.CurrentFilename())
	assert.Equal(t, line.Line{}, s.CurrentLine())
}

func TestFilenameStack(t *testing.T) {
	s, _ := newStatus(t)

	t.Run("PopFilename on empty stack is a no-op", func(t *testing.T) {
		s.PopFilename()
		assert.Empty(t, s.CurrentFilename())
	})

	t.Run("Push and Current", func(t *testing.T) {
		s.PushFilename("main.adv")
		assert.Equal(t, "main.adv", s.CurrentFilename())

		s.PushFilename("items.inc.adv")
		assert.Equal(t, "items.inc.adv", s.CurrentFilename())
	})

	t.Run("Pop returns to previous", func(t *testing.T) {
		s.PopFilename()
		assert.Equal(t, "main.adv", s.CurrentFilename())

		s.PopFilename()
		assert.Empty(t, s.CurrentFilename())
	})
}

func TestLineStackOverflow(t *testing.T) {
	s, _ := newStatus(t)

	// Push more than stackSize (5) lines; the oldest are evicted from the front.
	for i := 1; i <= 8; i++ {
		s.AppendStack(line.New("text", i))
	}

	require.Len(t, s.Stack, 5, "stack must be capped at 5 entries")
	assert.Equal(t, 4, s.Stack[0].Number(), "oldest retained line should be #4")
	assert.Equal(t, 8, s.Stack[len(s.Stack)-1].Number())
	assert.Equal(t, 8, s.CurrentLine().Number(), "CurrentLine returns the most recent")
}

func TestCommentAndMultiLine(t *testing.T) {
	s, _ := newStatus(t)

	t.Run("Comment lifecycle", func(t *testing.T) {
		s.SetComment(line.New("// a note", 1))
		assert.True(t, s.Comment.IsOn())

		s.UnsetComment()
		assert.False(t, s.Comment.IsOn())
	})

	t.Run("MultiLine appends", func(t *testing.T) {
		s.AppendLine(line.New(`desc: """`, 1))
		s.AppendLine(line.New(`some text`, 2))
		s.AppendLine(line.New(`"""`, 3))

		assert.True(t, s.MultiLine.IsOn())
		assert.Equal(t, 3, s.MultiLine.Len())
	})
}

func TestCurrentStoreableLifecycle(t *testing.T) {
	s, db := newStatus(t)

	t.Run("no current returns OK on save and false on get", func(t *testing.T) {
		assert.False(t, s.HasCurrent())

		err := s.SaveCurrentStoreable()
		assert.True(t, err.IsOK())

		assert.Zero(t, s.GetCurrentLabelID())

		var i item.Item
		assert.False(t, s.GetCurrentStoreable(&i))
	})

	t.Run("set, get, save", func(t *testing.T) {
		labelID, err := db.CreateLabel("sword")
		require.NoError(t, err)
		nounID, err := db.CreateLabel("noun-sword")
		require.NoError(t, err)

		// Push a line so the stored entity remembers where it was declared.
		s.AppendStack(line.New("excalibur: sword shiny", 10))
		s.PushFilename("items.inc.adv")

		it := item.New()
		it.LabelID = labelID
		it.NounID = nounID
		it.Description = "a sword"

		s.SetCurrentStoreable(it)
		require.True(t, s.HasCurrent())
		assert.Equal(t, labelID, s.GetCurrentLabelID())

		var got *item.Item
		require.True(t, s.GetCurrentStoreable(&got))
		assert.Equal(t, "a sword", got.Description)

		cerr := s.SaveCurrentStoreable()
		require.True(t, cerr.IsOK())

		// After save the current is cleared.
		assert.False(t, s.HasCurrent())
		assert.Zero(t, s.GetCurrentLabelID())

		// And the entity is persisted.
		assert.Equal(t, 1, db.CountRecordsByKind(kind.Item))
	})

	t.Run("SetCurrentStoreable swaps existing storeable", func(t *testing.T) {
		s, db := newStatus(t)

		labelID1, _ := db.CreateLabel("first")
		labelID2, _ := db.CreateLabel("second")
		nounID, _ := db.CreateLabel("noun")

		first := item.New()
		first.LabelID = labelID1
		first.NounID = nounID
		first.Description = "first"
		s.SetCurrentStoreable(first)

		second := item.New()
		second.LabelID = labelID2
		second.NounID = nounID
		second.Description = "second"
		s.SetCurrentStoreable(second)

		assert.Equal(t, labelID2, s.GetCurrentLabelID())
	})

	t.Run("GetCurrentStoreable with non-pointer dst returns false", func(t *testing.T) {
		s, db := newStatus(t)

		labelID, _ := db.CreateLabel("x")
		nounID, _ := db.CreateLabel("noun")
		it := item.New()
		it.LabelID = labelID
		it.NounID = nounID
		s.SetCurrentStoreable(it)

		var notAPointer item.Item
		assert.False(t, s.GetCurrentStoreable(notAPointer))
	})

	t.Run("ClearCurrent drops the entity without saving", func(t *testing.T) {
		s, db := newStatus(t)

		labelID, _ := db.CreateLabel("tmp")
		nounID, _ := db.CreateLabel("noun")
		it := item.New()
		it.LabelID = labelID
		it.NounID = nounID
		s.SetCurrentStoreable(it)
		require.True(t, s.HasCurrent())

		s.ClearCurrent()
		assert.False(t, s.HasCurrent())
		assert.Equal(t, 0, db.CountRecordsByKind(kind.Item),
			"ClearCurrent must not persist anything")
	})

	t.Run("SaveCurrentStoreable returns compiler error on db failure", func(t *testing.T) {
		s, _ := newStatus(t)

		// An item with a non-existing LabelID makes db.Create fail.
		it := item.New()
		it.LabelID = 9999
		it.NounID = 9999
		s.SetCurrentStoreable(it)

		cerr := s.SaveCurrentStoreable()
		assert.False(t, cerr.IsOK())
		// The current entity stays set so the caller can still inspect it.
		assert.True(t, s.HasCurrent())
	})
}

func TestCurrentPath(t *testing.T) {
	s, _ := newStatus(t)

	t.Run("with no filename joins with empty base", func(t *testing.T) {
		assert.Equal(t, "gfx/logo.png", s.CurrentPath("gfx/logo.png"))
	})

	t.Run("joins relative to current file dir", func(t *testing.T) {
		s.PushFilename(filepath.Join("adventures", "main.adv"))
		assert.Equal(t, filepath.Join("adventures", "gfx", "logo.png"),
			s.CurrentPath("gfx", "logo.png"))
	})

	t.Run("nested includes use the topmost filename", func(t *testing.T) {
		s.PushFilename(filepath.Join("adventures", "shared", "items.inc.adv"))
		assert.Equal(t, filepath.Join("adventures", "shared", "x.png"),
			s.CurrentPath("x.png"))
	})
}

func TestRunnedValidators(t *testing.T) {
	s, _ := newStatus(t)

	assert.False(t, s.HasRunValidator(kind.Item))

	s.FlagValidator(kind.Item)
	assert.True(t, s.HasRunValidator(kind.Item))

	// Flagging twice is idempotent (no duplicates).
	s.FlagValidator(kind.Item)
	s.FlagValidator(kind.Item)
	assert.True(t, s.HasRunValidator(kind.Item))

	// Other kinds are independent.
	assert.False(t, s.HasRunValidator(kind.Character))
	s.FlagValidator(kind.Character)
	assert.True(t, s.HasRunValidator(kind.Character))
}

func TestRunnedReplacers(t *testing.T) {
	s, _ := newStatus(t)

	assert.False(t, s.HasRunReplacer(kind.Word))

	s.FlagReplacer(kind.Word)
	assert.True(t, s.HasRunReplacer(kind.Word))

	s.FlagReplacer(kind.Word)
	assert.True(t, s.HasRunReplacer(kind.Word))

	assert.False(t, s.HasRunReplacer(kind.Item))
}
