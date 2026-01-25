package parser_test

import (
	"strings"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/jorgefuertes/thenewquill/internal/lang"
	"github.com/jorgefuertes/thenewquill/internal/parser"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	type phrase struct {
		str      string
		original string
		sub      bool
	}

	type testCase struct {
		name     string
		lang     lang.Lang
		input    string
		expected []phrase
	}

	testCases := []testCase{
		{
			"Spanish single word connection",
			lang.ES,
			"Norte",
			[]phrase{{"ir norte", "Norte", false}},
		},
		{
			"Spanish composed 3 phrases",
			lang.ES,
			"Coger la llave y abrir la puerta, ¡Ir norte!",
			[]phrase{
				{"coger llave", "Coger la llave", false},
				{"abrir puerta", "abrir la puerta", false},
				{"ir norte", "Ir norte", false},
			},
		},
		{
			"Spanish talking to character",
			lang.ES,
			`decir al elfo "dame el cuchillo y vete a tu casa"`,
			[]phrase{
				{"decir elfo", "decir al elfo", false},
				{"dame cuchillo", "dame el cuchillo", true},
				{"vete casa", "vete a tu casa", true},
			},
		},
		{
			"Spanish talking to character with infered verb",
			lang.ES,
			`elfo "dame el cuchillo y vete a tu casa"`,
			[]phrase{
				{"decir elfo", "elfo", false},
				{"dame cuchillo", "dame el cuchillo", true},
				{"vete casa", "vete a tu casa", true},
			},
		},
		{
			"Spanish exmine an object with infered verb",
			lang.ES,
			`cuchillo`,
			[]phrase{
				{"examinar cuchillo", "cuchillo", false},
			},
		},
		{
			"English single word connection",
			lang.EN,
			"North",
			[]phrase{{"go north", "North", false}},
		},
		{
			"English composed 3 phrases",
			lang.EN,
			"take the key and open the door. then go north!",
			[]phrase{
				{"take key", "take the key", false},
				{"open door", "open the door", false},
				{"go north", "go north", false},
			},
		},
		{
			"English talking to character",
			lang.EN,
			`say to Elf "give me the knife and go home"`,
			[]phrase{
				{"say elf", "say to Elf", false},
				{"give knife", "give me the knife", true},
				{"go home", "go home", true},
			},
		},
		{
			"English talking to character with infered verb",
			lang.EN,
			`Elf "give me the knife and go home"`,
			[]phrase{
				{"say elf", "Elf", false},
				{"give knife", "give me the knife", true},
				{"go home", "go home", true},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			p := setupParser(t, c.lang)
			p.Parse(c.input)

			var phrases []string
			for _, s := range p.Sentences {
				phrases = append(phrases, s.String())
			}

			var originals []string
			for _, s := range p.Sentences {
				originals = append(originals, s.Original())
			}

			require.Equal(t, len(c.expected), len(p.Sentences), "phrases: \n\t- %s\noriginals:\n\t- %s",
				strings.Join(phrases, "\n\t- "), strings.Join(originals, "\n\t- "))

			for i, e := range c.expected {
				lsStr := p.NextLS().String()
				require.Equal(t, e.str, lsStr, "phrase %d does not match", i+1)
				require.Equal(t, e.original, p.Current().Original(), "phrase %d original does not match", i+1)
				require.Equal(t, e.sub, p.Current().IsSub(), "phrase %d SUB flag does not match", i+1)
			}
		})
	}
}

func setupParser(t *testing.T, l lang.Lang) *parser.Parser {
	t.Helper()

	db := database.NewDB()
	configStore := config.NewService(db)
	configStore.Set(config.LanguageParamLabel, l.String())
	wordStore := word.NewService(db, configStore)

	if l == lang.EN {
		createWords(t, wordStore, word.Conjunction, "and", "then", "also")
		createWords(t, wordStore, word.Verb, "take", "open", "go", "give", "say", "examine")
		createWords(t, wordStore, word.Noun, "key", "door", "north", "elf", "home", "knife")
		setAsItem(t, wordStore, "key")
		setAsItem(t, wordStore, "knife")
		setAsItem(t, wordStore, "door")
		setAsConnection(t, wordStore, "north")
		setAsCharacter(t, wordStore, "elf")
	} else {
		createWords(t, wordStore, word.Conjunction, "y", "luego", "también", "después", "entonces")
		createWords(t, wordStore, word.Verb, "coger", "abrir", "ir", "dame", "vete", "decir", "examinar")
		createWords(t, wordStore, word.Noun, "llave", "puerta", "norte", "elfo", "casa", "cuchillo")
		setAsItem(t, wordStore, "llave")
		setAsItem(t, wordStore, "cuchillo")
		setAsItem(t, wordStore, "puerta")
		setAsConnection(t, wordStore, "norte")
		setAsCharacter(t, wordStore, "elfo")
	}

	p, err := parser.New(wordStore)
	require.NoError(t, err)
	require.NotNil(t, p)

	verr := wordStore.ValidateAll()
	for _, e := range verr {
		t.Error(e)
	}

	return p
}

func createWords(t *testing.T, wordStore *word.Service, wType word.WordType, labels ...string) {
	t.Helper()

	for _, s := range labels {
		labelID, err := wordStore.DB().CreateLabel(s)
		require.NoError(t, err)
		require.NotZero(t, labelID)

		w := word.New(labelID, wType, s)
		id, err := wordStore.Create(w)
		require.NoError(t, err)
		require.NotZero(t, id)
	}
}

func setAsItem(t *testing.T, wordStore *word.Service, label string) {
	t.Helper()

	w, err := wordStore.Get().WithLabel(label).First()
	require.NoError(t, err)
	require.NotNil(t, w)

	w.IsItem = true
	err = wordStore.Update(w)
	require.NoError(t, err)
}

func setAsConnection(t *testing.T, wordStore *word.Service, label string) {
	t.Helper()

	w, err := wordStore.Get().WithLabel(label).First()
	require.NoError(t, err)
	require.NotNil(t, w)

	w.IsConnection = true
	err = wordStore.Update(w)
	require.NoError(t, err)
}

func setAsCharacter(t *testing.T, wordStore *word.Service, label string) {
	t.Helper()

	w, err := wordStore.Get().WithLabel(label).First()
	require.NoError(t, err)
	require.NotNil(t, w)

	w.IsCharacter = true
	err = wordStore.Update(w)
	require.NoError(t, err)
}
