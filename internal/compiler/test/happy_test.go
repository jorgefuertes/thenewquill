package compiler_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"github.com/jorgefuertes/thenewquill/internal/compiler"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompilerHappyPath(t *testing.T) {
	a, err := compiler.Compile("src/happy/test.adv")
	require.NoError(t, err)

	t.Run("vars", func(t *testing.T) {
		// vars
		assert.Equal(t, 7, a.Variables.Count())

		testCases := []struct {
			key      string
			expected any
		}{
			{"testTrue", true},
			{"testFalse", false},
			{"number", 10},
			{"number2", 20},
			{"aFloat", 1.5},
			{"name", `The New Quill Adventure Writing System`},
			{"hello", `Hello, _.\nWelcome to _.\n`},
		}

		for _, tc := range testCases {
			t.Run(tc.key, func(t *testing.T) {
				actual, err := a.Variables.FindByLabel(tc.key)
				require.NoError(t, err)
				assert.EqualValues(t, tc.expected, actual)
			})
		}
	})

	t.Run("Vocabulary", func(t *testing.T) {
		t.Run("Not Found", func(t *testing.T) {
			_, err := a.Words.FindByLabel("foo")
			require.Error(t, err)
		})

		testCases := []struct {
			kind  word.WordType
			label string
			syns  []string
		}{
			{word.Verb, "_", []string{}},
			{word.Noun, "_", []string{}},
			{word.Pronoun, "_", []string{}},
			{word.Adjective, "_", []string{}},
			{word.Adverb, "_", []string{}},
			{word.Conjunction, "_", []string{}},
			{word.Preposition, "_", []string{}},
			{word.Verb, "entrar", []string{"entro", "entra"}},
			{word.Verb, "subir", []string{"sube", "subo"}},
			{word.Verb, "bajar", []string{"bajo", "baja"}},
			{word.Verb, "salir", []string{"salgo", "sal"}},
			{word.Verb, "coger", []string{"coge", "llevar", "recoger", "recoge", "pillar"}},
			{word.Verb, "dejar", []string{"dejo", "soltar", "suelta", "suelto"}},
			{word.Verb, "quitar", []string{"quito", "quita"}},
			{word.Verb, "poner", []string{"pongo", "pone"}},
			{word.Verb, "abrir", []string{"abro", "abre"}},
			{word.Verb, "cerrar", []string{"cierro", "cierra"}},
			{word.Verb, "fin", []string{"terminar", "acabar", "sistema"}},
			{word.Verb, "quill", []string{"github.com/jorgefuertes/thenewquill", "tnq"}},
			{word.Verb, "ad", []string{}},
			{word.Verb, "decir", []string{"di", "hablar", "habla"}},
			{word.Verb, "ir", []string{"voy", "vamos", "ve"}},
			{word.Verb, "ex", []string{"exam", "examinar", "examina", "mirar", "mira", "miro"}},
			{word.Verb, "save", []string{"grabar", "graba", "grabo", "salvar", "salva", "salvo"}},
			{word.Verb, "ram", []string{"ramsave"}},
			{word.Noun, "jugador", []string{}},
			{word.Noun, "enano", []string{}},
			{word.Noun, "elfo", []string{}},
			{word.Noun, "norte", []string{"n", "adelante"}},
			{word.Noun, "sur", []string{"s", "atrás"}},
			{word.Noun, "este", []string{"e"}},
			{word.Noun, "oeste", []string{"o", "w"}},
			{word.Noun, "abrigo", []string{"chaqueta"}},
			{word.Noun, "guantes", []string{}},
			{word.Noun, "cofre", []string{"arcón"}},
			{word.Noun, "monedas", []string{"dinero"}},
			{word.Noun, "denario", []string{}},
			{word.Noun, "carta", []string{"papel"}},
			{word.Noun, "antorcha", []string{"tea", "linterna"}},
			{word.Noun, "llave", []string{}},
			{word.Noun, "cinturón", []string{}},
			{word.Noun, "petaca", []string{"botella"}},
			{word.Noun, "talismán", []string{"amuleto"}},
			{word.Noun, "ropa", []string{"vestido"}},
			{word.Noun, "bolsa", []string{"saco", "petate"}},
			{word.Adjective, "encendida", []string{}},
			{word.Adjective, "apagada", []string{}},
			{word.Adjective, "dorada", []string{}},
			{word.Adjective, "plateada", []string{}},
			{word.Pronoun, "el", []string{"la", "los", "las"}},
			{word.Preposition, "dentro", []string{"adentro"}},
			{word.Conjunction, "enton", []string{"luego", "tras", "y"}},
		}

		for _, tc := range testCases {
			t.Run(tc.label, func(t *testing.T) {
				w, err := a.Words.FindByLabel(tc.label)
				require.NoError(t, err)
				require.NotNil(t, w)
				assert.Equal(t, tc.kind, w.Type)
				assert.Equal(t, tc.syns, w.Synonyms, "synonyms for %s doesn't match", tc.label)
				for _, syn := range tc.syns {
					wFromSyn, err := a.Words.First(syn)
					require.NoError(t, err)
					assert.Equal(t, w, wFromSyn)
					assert.True(t, w.Is(w.Type, syn))
				}
			})
		}

		for _, w := range a.Words.All() {
			existsHere := false
			for _, tc := range testCases {
				label, err := a.DB.GetLabel(w.ID)
				require.NoError(t, err)

				if label.ID == w.ID && tc.kind == w.Type {
					existsHere = true
					break
				}
			}

			assert.True(t, existsHere, "%s with ID %s doesn't exist in the test cases", w.Type.String(), w.ID)
		}
	})

	t.Run("Messages", func(t *testing.T) {
		testCases := []struct {
			label        string
			expected     string
			shouldExists bool
		}{
			{"dark", "No se vé nada.", true},
			{"loc-objects", "Aquí hay", true},
			{"not-needed", "No es necesario para jugar la aventura.", true},
			{"cant", "No puedes _", true},
			{"cant-do", "No puedes hacerlo.", true},
			{"err-save", "Error al grabar el fichero _.", true},
			{"filename", "Nombre del fichero:", true},
			{"foo", "", false},
			{"test", "This is a test message.", true},
			{"test2", `This is another \"test\" message.`, true},
			{
				"multiline",
				"This is a message with a heredoc string.\nLine 2.\nLine 3.\nLine 4 without carrige return at the " +
					"end.\n\tLine 5 and this line is indented.",
				true,
			},
			{"bar", "", false},
		}

		for _, tc := range testCases {
			t.Run(tc.label, func(t *testing.T) {
				m, err := a.Messages.FindByLabel(tc.label)
				if tc.shouldExists {
					require.NoError(t, err)
					require.NotNil(t, m)
					assert.Equal(t, tc.expected, m.Text)
				} else {
					require.Error(t, err)
				}
			})
		}
	})

	t.Run("Characters", func(t *testing.T) {
		testCases := []struct {
			label         string
			noun          string
			adj           string
			desc          string
			created       bool
			human         bool
			locationLabel string
			vars          map[string]any
		}{
			{"player", "jugador", "_", "Eres un pedazo de jugador", true, true, "celda", map[string]any{"health": 100}},
			{
				"enano",
				"enano",
				"_",
				"Un enano cabreado",
				true,
				false,
				"cuartelillo",
				map[string]any{"patience": 255, "death": false},
			},
			{"elfo", "elfo", "_", "Un elfo llorica", false, false, "gradas", map[string]any{"hidden": true}},
		}

		for _, tc := range testCases {
			t.Run(tc.label, func(t *testing.T) {
				c, err := a.Characters.FindByLabel(tc.label)
				require.NoError(t, err)
				require.NotNil(t, c)

				assert.Equal(t, tc.desc, c.Description)

				noun, err := a.Words.FindByLabel(tc.noun)
				require.NoError(t, err)
				assert.Equal(t, c.NounID, noun.ID)
				assert.Equal(t, word.Noun, noun.Type)

				adj, err := a.Words.FindByLabel(tc.adj)
				require.NoError(t, err)
				assert.Equal(t, c.AdjectiveID, adj.ID)
				assert.Equal(t, word.Adjective, adj.Type)

				loc, err := a.Locations.FindByLabel(tc.locationLabel)
				require.NoError(t, err)
				assert.Equal(t, loc.ID, c.LocationID)

				assert.Equal(t, tc.created, c.Created)
				assert.Equal(t, tc.human, c.Human)

				for k, v := range tc.vars {
					actual, err := a.Variables.FindByLabel(tc.label, k)
					require.NoError(t, err)
					assert.Equal(t, v, actual.Value)
				}
			})
		}
	})
}
