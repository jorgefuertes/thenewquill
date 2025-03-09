package compiler_test

import (
	"testing"

	"thenewquill/internal/adventure/words"
	"thenewquill/internal/compiler"

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
			{words.Verb, "quill", []string{"thenewquill", "tnq"}},
			{words.Verb, "ad", []string{}},
			{words.Verb, "decir", []string{"di", "hablar", "habla"}},
			{words.Verb, "ir", []string{"voy", "vamos", "ve"}},
			{words.Verb, "ex", []string{"exam", "examinar", "examina", "mirar", "mira", "miro"}},
			{words.Verb, "save", []string{"grabar", "graba", "grabo", "salvar", "salva", "salvo"}},
			{words.Verb, "ram", []string{"ramsave"}},
			{words.Noun, "norte", []string{"n", "adelante"}},
			{words.Noun, "sur", []string{"s", "atrás"}},
			{words.Noun, "este", []string{"e"}},
			{words.Noun, "oeste", []string{"o", "w"}},
			{words.Noun, "abrigo", []string{"chaqueta"}},
			{words.Noun, "guantes", []string{}},
			{words.Noun, "cofre", []string{"arcón"}},
			{words.Noun, "monedas", []string{"dinero"}},
			{words.Noun, "carta", []string{"papel"}},
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
				assert.Len(t, w.Synonyms, len(tc.syns))
				assert.Equal(t, tc.syns, w.Synonyms)
			})
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
}
