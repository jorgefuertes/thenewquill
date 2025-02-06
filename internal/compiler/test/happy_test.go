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
	a, err := compiler.Compile("adv_files/happy/test.adv")
	require.NoError(t, err)

	// vars
	assert.Equal(t, 7, a.Vars.Count())

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
				assert.Equal(t, tc.syns, w.Synonyms)
			})
		}
	})

	t.Run("Messages", func(t *testing.T) {
		testCases := []struct {
			kind          msg.MsgType
			label         string
			expected      string
			expectedError error
		}{
			{msg.SystemMsg, "dark", "No se vé nada.", nil},
			{msg.SystemMsg, "loc-objects", "Aquí hay", nil},
			{msg.SystemMsg, "not-needed", "No es necesario para jugar la aventura.", nil},
			{msg.SystemMsg, "cant", "No puedes _", nil},
			{msg.SystemMsg, "cant-do", "No puedes hacerlo.", nil},
			{msg.SystemMsg, "err-save", "Error al grabar el fichero _.", nil},
			{msg.SystemMsg, "filename", "Nombre del fichero:", nil},
			{msg.UserMsg, "foo", "", msg.ErrMsgNotFound},
			{msg.UserMsg, "test", "This is a test message.", nil},
			{msg.UserMsg, "test2", `This is another \"test\" message.`, nil},
			{
				msg.UserMsg,
				"multiline",
				"This is a message with a heredoc string.\nLine 2.\nLine 3.\nLine 4 without carrige return at the " +
					"end.\n\tLine 5 and this line is indented.",
				nil,
			},
			{msg.UserMsg, "bar", "", msg.ErrMsgNotFound},
		}

		for _, tc := range testCases {
			t.Run(tc.label, func(t *testing.T) {
				m, err := a.Messages.GetText(tc.kind, tc.label)
				if tc.expectedError != nil {
					require.Error(t, err)
					require.ErrorIs(t, err, tc.expectedError)

					return
				}

				require.NoError(t, err)
				assert.Equal(t, tc.expected, m)
			})
		}
	})
}
