package compiler_test

import (
	"testing"

	"thenewquill/internal/adventure/msg"
	"thenewquill/internal/adventure/voc"
	"thenewquill/internal/compiler"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompilerHappyPath(t *testing.T) {
	a, err := compiler.Compile("src/happy/test.adv")
	require.NoError(t, err)

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

	t.Run("vars", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.key, func(t *testing.T) {
				actual := a.Vars.Get(tc.key)
				assert.EqualValues(t, tc.expected, actual)
			})
		}
	})

	t.Run("Vocabulary", func(t *testing.T) {
		t.Run("Nil word", func(t *testing.T) {
			require.Nil(t, a.Vocabulary.Get(voc.Verb, "foo"))
		})

		testCases := []struct {
			kind  voc.WordType
			label string
			syns  []string
		}{
			{voc.Verb, "subir", []string{"sube", "subo"}},
			{voc.Verb, "bajar", []string{"bajo", "baja"}},
			{voc.Verb, "salir", []string{"salgo", "sal"}},
			{voc.Verb, "coger", []string{"coge", "llevar", "recoger", "recoge", "pillar"}},
			{voc.Verb, "dejar", []string{"dejo", "soltar", "suelta", "suelto"}},
			{voc.Verb, "quitar", []string{"quito", "quita"}},
			{voc.Verb, "poner", []string{"pongo", "pone"}},
			{voc.Verb, "abrir", []string{"abro", "abre"}},
			{voc.Verb, "cerrar", []string{"cierro", "cierra"}},
			{voc.Verb, "fin", []string{"terminar", "acabar", "sistema"}},
			{voc.Verb, "quill", []string{"thenewquill", "tnq"}},
			{voc.Verb, "ad", []string{}},
			{voc.Verb, "decir", []string{"di", "hablar", "habla"}},
			{voc.Verb, "ir", []string{"voy", "vamos", "ve"}},
			{voc.Verb, "ex", []string{"exam", "examinar", "examina", "mirar", "mira", "miro"}},
			{voc.Verb, "save", []string{"grabar", "graba", "grabo", "salvar", "salva", "salvo"}},
			{voc.Verb, "ram", []string{"ramsave"}},
			{voc.Noun, "norte", []string{"n", "adelante"}},
			{voc.Noun, "sur", []string{"s", "atrás"}},
			{voc.Noun, "este", []string{"e"}},
			{voc.Noun, "oeste", []string{"o", "w"}},
			{voc.Noun, "abrigo", []string{"chaqueta"}},
			{voc.Noun, "guantes", []string{}},
			{voc.Noun, "cofre", []string{"arcón"}},
			{voc.Noun, "monedas", []string{"dinero"}},
			{voc.Noun, "carta", []string{"papel"}},
			{voc.Pronoun, "el", []string{"la", "los", "las"}},
			{voc.Preposition, "dentro", []string{"adentro"}},
			{voc.Conjunction, "enton", []string{"luego", "tras", "y"}},
		}

		for _, tc := range testCases {
			t.Run(tc.label, func(t *testing.T) {
				w := a.Vocabulary.Get(tc.kind, tc.label)
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
			kind         msg.MsgType
			label        string
			expected     string
			shouldExists bool
		}{
			{msg.SystemMsg, "dark", "No se vé nada.", true},
			{msg.SystemMsg, "loc-objects", "Aquí hay", true},
			{msg.SystemMsg, "not-needed", "No es necesario para jugar la aventura.", true},
			{msg.SystemMsg, "cant", "No puedes _", true},
			{msg.SystemMsg, "cant-do", "No puedes hacerlo.", true},
			{msg.SystemMsg, "err-save", "Error al grabar el fichero _.", true},
			{msg.SystemMsg, "filename", "Nombre del fichero:", true},
			{msg.UserMsg, "foo", "", false},
			{msg.UserMsg, "test", "This is a test message.", true},
			{msg.UserMsg, "test2", `This is another \"test\" message.`, true},
			{
				msg.UserMsg,
				"multiline",
				"This is a message with a heredoc string.\nLine 2.\nLine 3.\nLine 4 without carrige return at the " +
					"end.\n\tLine 5 and this line is indented.",
				true,
			},
			{msg.UserMsg, "bar", "", false},
		}

		for _, tc := range testCases {
			t.Run(tc.label, func(t *testing.T) {
				m := a.Messages.Get(tc.kind, tc.label)
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
