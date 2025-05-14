package compiler_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/words"
	"github.com/jorgefuertes/thenewquill/internal/compiler"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompilerHappyPath(t *testing.T) {
	a, err := compiler.Compile("src/happy/test.adv")
	require.NoError(t, err)

	t.Run("vars", func(t *testing.T) {
		// vars
		assert.Equal(t, 7, a.Vars.Len())

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
				actual := a.Vars.Get(tc.key)
				assert.EqualValues(t, tc.expected, actual)
			})
		}
	})

	t.Run("Vocabulary", func(t *testing.T) {
		t.Run("Nil word", func(t *testing.T) {
			require.Nil(t, a.Words.Get(words.Verb, "foo"))
		})

		testCases := []struct {
			kind  words.WordType
			label string
			syns  []string
		}{
			{words.Unknown, "_", []string{}},
			{words.Verb, "_", []string{}},
			{words.Noun, "_", []string{}},
			{words.Pronoun, "_", []string{}},
			{words.Adjective, "_", []string{}},
			{words.Adverb, "_", []string{}},
			{words.Conjunction, "_", []string{}},
			{words.Preposition, "_", []string{}},
			{words.Verb, "entrar", []string{"entro", "entra"}},
			{words.Verb, "subir", []string{"sube", "subo"}},
			{words.Verb, "bajar", []string{"bajo", "baja"}},
			{words.Verb, "salir", []string{"salgo", "sal"}},
			{words.Verb, "coger", []string{"coge", "llevar", "recoger", "recoge", "pillar"}},
			{words.Verb, "dejar", []string{"dejo", "soltar", "suelta", "suelto"}},
			{words.Verb, "quitar", []string{"quito", "quita"}},
			{words.Verb, "poner", []string{"pongo", "pone"}},
			{words.Verb, "abrir", []string{"abro", "abre"}},
			{words.Verb, "cerrar", []string{"cierro", "cierra"}},
			{words.Verb, "fin", []string{"terminar", "acabar", "sistema"}},
			{words.Verb, "quill", []string{"github.com/jorgefuertes/thenewquill", "tnq"}},
			{words.Verb, "ad", []string{}},
			{words.Verb, "decir", []string{"di", "hablar", "habla"}},
			{words.Verb, "ir", []string{"voy", "vamos", "ve"}},
			{words.Verb, "ex", []string{"exam", "examinar", "examina", "mirar", "mira", "miro"}},
			{words.Verb, "save", []string{"grabar", "graba", "grabo", "salvar", "salva", "salvo"}},
			{words.Verb, "ram", []string{"ramsave"}},
			{words.Noun, "jugador", []string{}},
			{words.Noun, "enano", []string{}},
			{words.Noun, "elfo", []string{}},
			{words.Noun, "norte", []string{"n", "adelante"}},
			{words.Noun, "sur", []string{"s", "atrás"}},
			{words.Noun, "este", []string{"e"}},
			{words.Noun, "oeste", []string{"o", "w"}},
			{words.Noun, "abrigo", []string{"chaqueta"}},
			{words.Noun, "guantes", []string{}},
			{words.Noun, "cofre", []string{"arcón"}},
			{words.Noun, "monedas", []string{"dinero"}},
			{words.Noun, "denario", []string{}},
			{words.Noun, "carta", []string{"papel"}},
			{words.Noun, "antorcha", []string{"tea", "linterna"}},
			{words.Noun, "llave", []string{}},
			{words.Noun, "cinturón", []string{}},
			{words.Noun, "petaca", []string{"botella"}},
			{words.Noun, "talismán", []string{"amuleto"}},
			{words.Noun, "ropa", []string{"vestido"}},
			{words.Noun, "bolsa", []string{"saco", "petate"}},
			{words.Adjective, "encendida", []string{}},
			{words.Adjective, "apagada", []string{}},
			{words.Adjective, "dorada", []string{}},
			{words.Adjective, "plateada", []string{}},
			{words.Pronoun, "el", []string{"la", "los", "las"}},
			{words.Preposition, "dentro", []string{"adentro"}},
			{words.Conjunction, "enton", []string{"luego", "tras", "y"}},
		}

		for _, tc := range testCases {
			t.Run(tc.label, func(t *testing.T) {
				w := a.Words.Get(tc.kind, tc.label)
				require.NotNil(t, w)
				assert.Equal(t, tc.label, w.Label)
				assert.Equal(t, tc.kind, w.Type)
				assert.Equal(t, tc.syns, w.Synonyms, "synonyms for %s doesn't match", tc.label)
				for _, syn := range tc.syns {
					wFromSyn := a.Words.First(syn)
					assert.Equal(t, w, wFromSyn)
					assert.True(t, w.Is(syn))
				}
			})
		}

		for _, w := range a.Words {
			existsHere := false
			for _, tc := range testCases {
				if tc.label == w.Label && tc.kind == w.Type {
					existsHere = true
					break
				}
			}

			assert.True(t, existsHere, "%s with label %s doesn't exist in the test cases", w.Type.String(), w.Label)
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
				m := a.Messages.Get(tc.label)
				if tc.shouldExists {
					require.NotNil(t, m)
					assert.Equal(t, tc.expected, m.Text)
				} else {
					require.Nil(t, m)
				}
			})
		}
	})

	t.Run("Characters", func(t *testing.T) {
		testCases := []struct {
			label         string
			nameLabel     string
			adjLabel      string
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
				c := a.Chars.Get(tc.label)
				require.NotNil(t, c)
				assert.Equal(t, tc.label, c.Label)
				assert.Equal(t, tc.nameLabel, c.Name.Label)
				assert.Equal(t, tc.adjLabel, c.Adjective.Label)
				assert.Equal(t, tc.desc, c.Description)
				assert.Equal(t, tc.locationLabel, c.Location.Label)
				assert.Equal(t, tc.created, c.Created)
				assert.Equal(t, tc.human, c.Human)
				assert.Equal(t, tc.vars, c.Vars.GetAll())
			})
		}
	})
}
