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
		str    string
		isMain bool
	}

	type testCase struct {
		name     string
		lang     lang.Lang
		input    string
		expected []phrase
	}

	testCases := []testCase{
		{
			"Spanish composed 3 phrases",
			lang.ES,
			"coger la llave y abrir la puerta. luego ¡Ir norte!",
			[]phrase{
				{"coger llave", true},
				{"abrir puerta", true},
				{"ir norte", true},
			},
		},
		{
			"Spanish single word connection",
			lang.ES,
			"Norte",
			[]phrase{{"norte", true}},
		},
		{
			"Spanish talking to character",
			lang.ES,
			"decir al elfo \"dame el cuchillo y vete a tu casa\"",
			[]phrase{
				{"hablar elfo", true},
				{"dame cuchillo", true},
				{"vete casa", true},
			},
		},
		{
			"English single word connection",
			lang.EN,
			"North",
			[]phrase{{"north", true}},
		},
		{
			"English composed 3 phrases",
			lang.EN,
			"take the key and open the door. then go north!",
			[]phrase{
				{"take key", true},
				{"open door", true},
				{"go north", true},
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

			require.Equal(t, len(c.expected), len(p.Sentences), "prhases: \n\t%s", strings.Join(phrases, "\n\t"))

			for i, e := range c.expected {
				lsStr := p.NextLS().String()
				require.Equal(t, e.str, lsStr, "phrase %d does not match", i+1)
				require.Equal(t, e.isMain, p.Current().IsMain(), "phrase %d main flag does not match", i+1)
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
		createWords(t, wordStore, word.Verb, "take", "open", "go", "give", "say")
		createWords(t, wordStore, word.Noun, "key", "door", "north", "elf", "home", "knife")
		setAsItem(t, wordStore, "key", "door", "knife")
		setAsConnection(t, wordStore, "north")
		setAsCharacter(t, wordStore, "elf")
	} else {
		createWords(t, wordStore, word.Conjunction, "y", "luego", "también", "después", "entonces")
		createWords(t, wordStore, word.Verb, "coger", "abrir", "ir", "dame", "vete", "decir")
		createWords(t, wordStore, word.Noun, "llave", "puerta", "norte", "elfo", "casa", "cuchillo")
		setAsItem(t, wordStore, "llave", "puerta", "cuchillo")
		setAsConnection(t, wordStore, "norte")
		setAsCharacter(t, wordStore, "elfo")
	}

	p, err := parser.New(wordStore)
	require.NoError(t, err)
	require.NotNil(t, p)

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

func setAsItem(t *testing.T, wordStore *word.Service, labels ...string) {
	t.Helper()

	for _, label := range labels {
		w, err := wordStore.Get().WithSynonym(label).First()
		require.NoError(t, err)
		require.NotNil(t, w)

		w.IsItem = true
		err = wordStore.Update(w)
		require.NoError(t, err)
	}
}

func setAsConnection(t *testing.T, wordStore *word.Service, label string) {
	t.Helper()

	w, err := wordStore.Get().WithSynonym(label).First()
	require.NoError(t, err)
	require.NotNil(t, w)

	w.IsConnection = true
	err = wordStore.Update(w)
	require.NoError(t, err)
}

func setAsCharacter(t *testing.T, wordStore *word.Service, label string) {
	t.Helper()

	w, err := wordStore.Get().WithSynonym(label).First()
	require.NoError(t, err)
	require.NotNil(t, w)

	w.IsCharacter = true
	err = wordStore.Update(w)
	require.NoError(t, err)
}
