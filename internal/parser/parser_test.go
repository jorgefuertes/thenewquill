package parser_test

import (
	"strings"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	"github.com/jorgefuertes/thenewquill/internal/adventure/item"
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
			"Spanish examine an object with infered verb",
			lang.ES,
			`cuchillo`,
			[]phrase{
				{"exami cuchillo", "cuchillo", false},
			},
		},
		// -- Spec examples: simple (docs/specs_v2.md lines 85-95) --
		{
			"Spec: coger hacha",
			lang.ES,
			"coger hacha",
			[]phrase{{"coger hacha", "coger hacha", false}},
		},
		{
			"Spec: examinar troll",
			lang.ES,
			"examinar troll",
			[]phrase{{"exami troll", "examinar troll", false}},
		},
		{
			"Spec: ex cueva (verb synonym)",
			lang.ES,
			"ex cueva",
			[]phrase{{"exami cueva", "ex cueva", false}},
		},
		// -- Spec examples: compound (docs/specs_v2.md lines 99-126) --
		{
			"Spec: SUBSLs with trailing SL",
			lang.ES,
			`decir hobbit "coge el hacha y corta leña" y coger leña`,
			[]phrase{
				{"decir hobbit", "decir hobbit", false},
				{"coger hacha", "coge el hacha", true},
				{"corta leña", "corta leña", true},
				{"coger leña", "y coger leña", false},
			},
		},
		{
			"Spec: adverb after noun",
			lang.ES,
			"coger la espada rápidamente y atacar al troll",
			[]phrase{
				{"coger espada rápidamente", "coger la espada rápidamente", false},
				{"ataca troll", "atacar al troll", false},
			},
		},
		{
			"Spec: adverb before noun same SL",
			lang.ES,
			"coger rápidamente la espada y atacar al troll",
			[]phrase{
				{"coger rápidamente espada", "coger rápidamente la espada", false},
				{"ataca troll", "atacar al troll", false},
			},
		},
		{
			"Spec: verb-only SL after conjunction",
			lang.ES,
			"sonreir al troll y bailar",
			[]phrase{
				{"sonre troll", "sonreir al troll", false},
				{"baila", "bailar", false},
			},
		},
		{
			"Spec: clitic truncation matarlo -> matar -> atacar",
			lang.ES,
			"coger rápidamente la espada y perseguir al troll y matarlo",
			[]phrase{
				{"coger rápidamente espada", "coger rápidamente la espada", false},
				{"perse troll", "perseguir al troll", false},
				{"ataca", "matarlo", false},
			},
		},
		// -- English tests --
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
	_, err := configStore.Set(config.LanguageParamLabel, l.String())
	require.NoError(t, err, "failed to set language config")
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
		createWordWithSynonyms(t, wordStore, word.Verb, "coger", "coger", "coge")
		createWordWithSynonyms(t, wordStore, word.Verb, "examinar", "examinar", "ex")
		createWordWithSynonyms(t, wordStore, word.Verb, "atacar", "atacar", "matar")
		createWords(t, wordStore, word.Verb, "abrir", "ir", "dame", "vete", "decir",
			"cortar", "sonreir", "bailar", "perseguir")
		createWords(t, wordStore, word.Noun, "llave", "puerta", "norte", "elfo", "casa",
			"cuchillo", "hacha", "troll", "cueva", "espada", "leña", "hobbit")
		createWords(t, wordStore, word.Adverb, "rápidamente")
		setAsItem(t, wordStore, "llave")
		setAsItem(t, wordStore, "cuchillo")
		setAsItem(t, wordStore, "puerta")
		setAsItem(t, wordStore, "hacha")
		setAsItem(t, wordStore, "espada")
		setAsItem(t, wordStore, "leña")
		setAsConnection(t, wordStore, "norte")
		setAsCharacter(t, wordStore, "elfo")
		setAsCharacter(t, wordStore, "troll")
		setAsCharacter(t, wordStore, "hobbit")
	}

	itemStore := item.NewService(db)
	charStore := character.NewService(db)

	p, err := parser.New(wordStore, itemStore, charStore)
	require.NoError(t, err)
	require.NotNil(t, p)

	verr := wordStore.ValidateAll()
	for _, e := range verr {
		t.Error(e)
	}

	return p
}

func createWordWithSynonyms(t *testing.T, wordStore *word.Service, wType word.WordType, label string, synonyms ...string) {
	t.Helper()

	labelID, err := wordStore.DB().CreateLabel(label)
	require.NoError(t, err)
	require.NotZero(t, labelID)

	w := word.New(labelID, wType, synonyms...)
	id, err := wordStore.Create(w)
	require.NoError(t, err)
	require.NotZero(t, id)
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

func TestBindings(t *testing.T) {
	db := database.NewDB()
	configStore := config.NewService(db)
	_, err := configStore.Set(config.LanguageParamLabel, lang.ES.String())
	require.NoError(t, err)
	wordStore := word.NewService(db, configStore)

	createWords(t, wordStore, word.Conjunction, "y")
	createWordWithSynonyms(t, wordStore, word.Verb, "coger", "coger", "coge")
	createWordWithSynonyms(t, wordStore, word.Verb, "examinar", "examinar", "ex")
	createWordWithSynonyms(t, wordStore, word.Verb, "atacar", "atacar", "matar")
	createWords(t, wordStore, word.Verb, "decir", "bailar")
	createWords(t, wordStore, word.Noun, "espada", "llave", "elfo", "troll")
	setAsItem(t, wordStore, "espada")
	setAsItem(t, wordStore, "llave")
	setAsCharacter(t, wordStore, "elfo")
	setAsCharacter(t, wordStore, "troll")

	itemStore := item.NewService(db)
	charStore := character.NewService(db)

	wordEspada, err := wordStore.Get().WithLabel("espada").First()
	require.NoError(t, err)
	wordLlave, err := wordStore.Get().WithLabel("llave").First()
	require.NoError(t, err)
	wordElfo, err := wordStore.Get().WithLabel("elfo").First()
	require.NoError(t, err)
	wordTroll, err := wordStore.Get().WithLabel("troll").First()
	require.NoError(t, err)

	espadaLabel, err := db.CreateLabel("item-espada")
	require.NoError(t, err)
	itemEspada := item.New()
	itemEspada.LabelID = espadaLabel
	itemEspada.NounID = wordEspada.ID
	itemEspada.Description = "a sword"
	_, err = itemStore.Create(itemEspada)
	require.NoError(t, err)

	llaveLabel, err := db.CreateLabel("item-llave")
	require.NoError(t, err)
	itemLlave := item.New()
	itemLlave.LabelID = llaveLabel
	itemLlave.NounID = wordLlave.ID
	itemLlave.Description = "a key"
	_, err = itemStore.Create(itemLlave)
	require.NoError(t, err)

	elfoLabel, err := db.CreateLabel("npc-elfo")
	require.NoError(t, err)
	npcElfo := character.New()
	npcElfo.LabelID = elfoLabel
	npcElfo.NounID = wordElfo.ID
	npcElfo.AdjectiveID = 1
	npcElfo.Description = "an elf"
	npcElfo.LocationID = 1
	_, err = db.Create(npcElfo)
	require.NoError(t, err)

	trollLabel, err := db.CreateLabel("npc-troll")
	require.NoError(t, err)
	npcTroll := character.New()
	npcTroll.LabelID = trollLabel
	npcTroll.NounID = wordTroll.ID
	npcTroll.AdjectiveID = 1
	npcTroll.Description = "a troll"
	npcTroll.LocationID = 1
	_, err = db.Create(npcTroll)
	require.NoError(t, err)

	newParser := func(t *testing.T) *parser.Parser {
		t.Helper()
		p, err := parser.New(wordStore, itemStore, charStore)
		require.NoError(t, err)
		return p
	}

	t.Run("item binding on single sentence", func(t *testing.T) {
		p := newParser(t)
		p.Parse("coger espada")

		require.Equal(t, 1, p.Len())
		ls := p.NextLS()
		require.NotNil(t, ls.Item, "Item should be bound")
		require.Equal(t, itemEspada.ID, ls.Item.ID)
		require.Nil(t, ls.NPC, "NPC should be nil")
	})

	t.Run("NPC binding on single sentence", func(t *testing.T) {
		p := newParser(t)
		p.Parse("decir elfo")

		require.Equal(t, 1, p.Len())
		ls := p.NextLS()
		require.NotNil(t, ls.NPC, "NPC should be bound")
		require.Equal(t, npcElfo.ID, ls.NPC.ID)
		require.Nil(t, ls.Item, "Item should be nil")
	})

	t.Run("item and NPC in separate sentences", func(t *testing.T) {
		p := newParser(t)
		p.Parse("coger espada y atacar troll")

		require.Equal(t, 2, p.Len())

		ls1 := p.NextLS()
		require.NotNil(t, ls1.Item, "first sentence should bind item")
		require.Equal(t, itemEspada.ID, ls1.Item.ID)
		require.Nil(t, ls1.NPC)

		ls2 := p.NextLS()
		require.NotNil(t, ls2.NPC, "second sentence should bind NPC")
		require.Equal(t, npcTroll.ID, ls2.NPC.ID)
		require.Nil(t, ls2.Item)
	})

	t.Run("item carry-forward to verb-only sentence", func(t *testing.T) {
		p := newParser(t)
		p.Parse("coger espada y bailar")

		require.Equal(t, 2, p.Len())

		ls1 := p.NextLS()
		require.NotNil(t, ls1.Item)
		require.Equal(t, itemEspada.ID, ls1.Item.ID)

		ls2 := p.NextLS()
		require.NotNil(t, ls2.Item, "item should carry forward")
		require.Equal(t, itemEspada.ID, ls2.Item.ID)
	})

	t.Run("NPC carry-forward to verb-only sentence", func(t *testing.T) {
		p := newParser(t)
		p.Parse("atacar troll y bailar")

		require.Equal(t, 2, p.Len())

		ls1 := p.NextLS()
		require.NotNil(t, ls1.NPC)
		require.Equal(t, npcTroll.ID, ls1.NPC.ID)

		ls2 := p.NextLS()
		require.NotNil(t, ls2.NPC, "NPC should carry forward")
		require.Equal(t, npcTroll.ID, ls2.NPC.ID)
	})

	t.Run("binding changes between sentences", func(t *testing.T) {
		p := newParser(t)
		p.Parse("coger espada y coger llave")

		require.Equal(t, 2, p.Len())

		ls1 := p.NextLS()
		require.NotNil(t, ls1.Item)
		require.Equal(t, itemEspada.ID, ls1.Item.ID)

		ls2 := p.NextLS()
		require.NotNil(t, ls2.Item)
		require.Equal(t, itemLlave.ID, ls2.Item.ID)
	})

	t.Run("no binding when noun is not an item or NPC", func(t *testing.T) {
		createWords(t, wordStore, word.Noun, "cueva")
		p := newParser(t)
		p.Parse("examinar cueva")

		require.Equal(t, 1, p.Len())
		ls := p.NextLS()
		require.Nil(t, ls.Item)
		require.Nil(t, ls.NPC)
	})
}
