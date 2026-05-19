package line_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/message"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsInclude(t *testing.T) {
	testCases := []struct {
		name string
		text string
		want string
		ok   bool
	}{
		{"valid include", `INCLUDE "items.inc.adv"`, "items.inc.adv", true},
		{"with subpath", `INCLUDE "shared/words.adv"`, "shared/words.adv", true},
		{"lowercase fails", `include "x.adv"`, "", false},
		{"no quotes fails", `INCLUDE x.adv`, "", false},
		{"garbage", `not an include`, "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := line.New(tc.text, 0).AsInclude()
			assert.Equal(t, tc.ok, ok)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAsSection(t *testing.T) {
	testCases := []struct {
		name string
		text string
		want kind.Kind
		ok   bool
	}{
		{"SECTION CONFIG", "SECTION CONFIG", kind.Param, true},
		{"section variables (lowercase)", "section variables", kind.Variable, true},
		{"SECTION ITEMS", "SECTION ITEMS", kind.Item, true},
		{"SECTION WORDS", "SECTION WORDS", kind.Word, true},
		{"SECTION MESSAGES", "SECTION MESSAGES", kind.Message, true},
		{"SECTION LOCATIONS", "SECTION LOCATIONS", kind.Location, true},
		{"SECTION CHARACTERS", "SECTION CHARACTERS", kind.Character, true},
		{"SECTION PICTURES", "SECTION PICTURES", kind.Blob, true},
		{"unknown section -> None match", "SECTION UNKNOWN", kind.None, true},
		{"non-section line", "title: \"x\"", kind.None, false},
		{"blank", "", kind.None, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := line.New(tc.text, 0).AsSection()
			assert.Equal(t, tc.ok, ok)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAsWord(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		wantKind string
		wantSyns []string
		ok       bool
	}{
		{
			name:     "verb with synonyms",
			text:     "verb: north, n",
			wantKind: "verb",
			wantSyns: []string{"north", "n"},
			ok:       true,
		},
		{
			name:     "single synonym",
			text:     "noun: sword",
			wantKind: "noun",
			wantSyns: []string{"sword"},
			ok:       true,
		},
		{
			name:     "trims extra spaces",
			text:     "adjective:   shiny,   bright",
			wantKind: "adjective",
			wantSyns: []string{"shiny", "bright"},
			ok:       true,
		},
		{name: "no colon", text: "verb north", ok: false},
		{name: "no synonyms", text: "verb:", ok: false},
		{name: "garbage", text: "??", ok: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			k, syns, ok := line.New(tc.text, 0).AsWord()
			assert.Equal(t, tc.ok, ok)

			if !ok {
				return
			}

			assert.Equal(t, tc.wantKind, k)
			assert.Equal(t, tc.wantSyns, syns)
		})
	}
}

func TestAsMsg(t *testing.T) {
	testCases := []struct {
		name       string
		text       string
		wantLabel  string
		wantText   string
		wantPlural message.Plural
		ok         bool
	}{
		{
			name:       "single line message",
			text:       `greeting: "Hello there"`,
			wantLabel:  "greeting",
			wantText:   "Hello there",
			wantPlural: message.Zero,
			ok:         true,
		},
		{
			name:       "plural one",
			text:       `count.one: "one item"`,
			wantLabel:  "count",
			wantText:   "one item",
			wantPlural: message.One,
			ok:         true,
		},
		{
			name:       "plural many",
			text:       `count.many: "many items"`,
			wantLabel:  "count",
			wantText:   "many items",
			wantPlural: message.Many,
			ok:         true,
		},
		{
			name:       "plural zero",
			text:       `count.zero: "no items"`,
			wantLabel:  "count",
			wantText:   "no items",
			wantPlural: message.Zero,
			ok:         true,
		},
		{name: "garbage", text: "not a message", ok: false},
		{name: "missing quotes", text: "greeting: hello", ok: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			label, text, plural, ok := line.New(tc.text, 0).AsMsg()
			assert.Equal(t, tc.ok, ok)

			if !ok {
				return
			}

			assert.Equal(t, tc.wantLabel, label)
			assert.Equal(t, tc.wantText, text)
			assert.Equal(t, tc.wantPlural, plural)
		})
	}
}

func TestAsLocationLabel(t *testing.T) {
	testCases := []struct {
		name string
		text string
		want string
		ok   bool
	}{
		{"plain label", "entrance:", "entrance", true},
		{"with leading whitespace", "  entrance:", "entrance", true},
		{"with trailing whitespace", "entrance:   ", "entrance", true},
		{"label with hyphen", "great-hall:", "great-hall", true},
		{"label with dot", "level1.entry:", "level1.entry", true},
		{"has trailing content -> not a label", `entrance: "something"`, "", false},
		{"garbage", "no colon here", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := line.New(tc.text, 0).AsLocationLabel()
			assert.Equal(t, tc.ok, ok)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAsLocationTitleAndDescription(t *testing.T) {
	t.Run("title", func(t *testing.T) {
		got, ok := line.New(`title: "The hall"`, 0).AsLocationTitle()
		require.True(t, ok)
		assert.Equal(t, "The hall", got)
	})

	t.Run("title with indentation", func(t *testing.T) {
		got, ok := line.New(`	title: "The hall"`, 0).AsLocationTitle()
		require.True(t, ok)
		assert.Equal(t, "The hall", got)
	})

	t.Run("description", func(t *testing.T) {
		got, ok := line.New(`desc: "A long hallway"`, 0).AsLocationDescription()
		require.True(t, ok)
		assert.Equal(t, "A long hallway", got)
	})

	t.Run("title miss", func(t *testing.T) {
		_, ok := line.New(`desc: "x"`, 0).AsLocationTitle()
		assert.False(t, ok)
	})
}

func TestAsLocationConns(t *testing.T) {
	t.Run("single exit", func(t *testing.T) {
		exits, ok := line.New("exits: north hall", 0).AsLocationConns()
		require.True(t, ok)
		assert.Equal(t, map[string]string{"north": "hall"}, exits)
	})

	t.Run("multiple exits", func(t *testing.T) {
		exits, ok := line.New("exits: north hall, south cave, east garden", 0).AsLocationConns()
		require.True(t, ok)
		assert.Equal(t, map[string]string{
			"north": "hall",
			"south": "cave",
			"east":  "garden",
		}, exits)
	})

	t.Run("conns alias", func(t *testing.T) {
		exits, ok := line.New("conns: up tower", 0).AsLocationConns()
		require.True(t, ok)
		assert.Equal(t, map[string]string{"up": "tower"}, exits)
	})

	t.Run("connections alias", func(t *testing.T) {
		exits, ok := line.New("connections: down dungeon", 0).AsLocationConns()
		require.True(t, ok)
		assert.Equal(t, map[string]string{"down": "dungeon"}, exits)
	})

	t.Run("garbage", func(t *testing.T) {
		exits, ok := line.New("not an exit line", 0).AsLocationConns()
		assert.False(t, ok)
		assert.Empty(t, exits)
	})
}

func TestAsLabelNounAdjDeclaration(t *testing.T) {
	testCases := []struct {
		name  string
		text  string
		wantL string
		wantN string
		wantA string
		ok    bool
	}{
		{
			name:  "valid declaration",
			text:  "excalibur: sword shiny",
			wantL: "excalibur", wantN: "sword", wantA: "shiny", ok: true,
		},
		{
			name:  "with indentation",
			text:  "   hero: man brave",
			wantL: "hero", wantN: "man", wantA: "brave", ok: true,
		},
		{name: "missing adjective", text: "excalibur: sword", ok: false},
		{name: "no colon", text: "excalibur sword shiny", ok: false},
		{name: "empty", text: "", ok: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l, n, a, ok := line.New(tc.text, 0).AsLabelNounAdjDeclaration()
			assert.Equal(t, tc.ok, ok)

			if !ok {
				return
			}

			assert.Equal(t, tc.wantL, l)
			assert.Equal(t, tc.wantN, n)
			assert.Equal(t, tc.wantA, a)
		})
	}
}

func TestAsConfig(t *testing.T) {
	testCases := []struct {
		name      string
		text      string
		wantLabel string
		wantValue string
		ok        bool
	}{
		{"title", `title: "Some title"`, "title", "Some title", true},
		{"author", `author: "Queru"`, "author", "Queru", true},
		{"description", `description: "A test adventure"`, "description", "A test adventure", true},
		{"version", `version: "1.0.0"`, "version", "1.0.0", true},
		{"date", `date: "2026-01-01"`, "date", "2026-01-01", true},
		{"language", `language: "es"`, "language", "es", true},
		{"unknown field", `nonsense: "x"`, "", "", false},
		{"garbage", `no colon`, "", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			label, value, ok := line.New(tc.text, 0).AsConfig()
			assert.Equal(t, tc.ok, ok)

			if !ok {
				return
			}

			assert.Equal(t, tc.wantLabel, label)
			assert.Equal(t, tc.wantValue, value)
		})
	}
}

func TestAsBlob(t *testing.T) {
	testCases := []struct {
		name      string
		text      string
		wantLabel string
		wantPath  string
		ok        bool
	}{
		{"valid blob", "logo: gfx/logo.png", "logo", "gfx/logo.png", true},
		{"with nested path", "tune: snd/intro/main.mp3", "tune", "snd/intro/main.mp3", true},
		{"with indentation", "  portrait: gfx/queru.jpg", "portrait", "gfx/queru.jpg", true},
		{"path without slash fails", "logo: logo.png", "", "", false},
		{"no extension fails", "logo: gfx/logo", "", "", false},
		{"no colon fails", "logo gfx/logo.png", "", "", false},
		{"garbage", "nope", "", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			label, path, ok := line.New(tc.text, 0).AsBlob()
			assert.Equal(t, tc.ok, ok)

			if !ok {
				return
			}

			assert.Equal(t, tc.wantLabel, label)
			assert.Equal(t, tc.wantPath, path)
		})
	}
}
