package compiler_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
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
		assert.Equal(t, 15, a.Variables.Count())

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
			{"gradas.people", 500},
			{"subasta.running", true},
			{"via.open", true},
			{"ant-on.on", false},
			{"player.health", 100},
			{"enano.patience", 255},
			{"enano.death", false},
			{"elfo.hidden", true},
		}

		for _, tc := range testCases {
			t.Run(tc.key, func(t *testing.T) {
				actual, err := a.Variables.FindByLabel(tc.key)
				require.NoError(t, err)

				switch reflect.TypeOf(tc.expected).Kind() {
				case reflect.Int:
					assert.Equal(t, tc.expected, actual.Int())
				case reflect.Float32, reflect.Float64:
					assert.Equal(t, tc.expected, actual.Float())
				case reflect.Bool:
					assert.Equal(t, tc.expected, actual.Bool())
				default:
					assert.EqualValues(t, tc.expected, actual.String())
				}
			})
		}

		vars := a.Variables.All()
		for _, v := range vars {
			exists := false

			for _, tc := range testCases {
				if a.DB.GetLabelName(v.ID) == tc.key {
					exists = true
				}
			}

			if !exists {
				t.Errorf("extra variable '%s = %v'", a.DB.GetLabelName(v.ID), v.Value)
			}
		}
	})

	t.Run("Vocabulary", func(t *testing.T) {
		t.Run("Not Found", func(t *testing.T) {
			_, err := a.Words.FindByLabel("foo")
			require.Error(t, err)
		})

		testCases := []struct {
			kind word.WordType
			syns []string
		}{
			{word.Verb, []string{"entrar", "entro", "entra"}},
			{word.Verb, []string{"subir", "sube", "subo"}},
			{word.Verb, []string{"bajar", "bajo", "baja"}},
			{word.Verb, []string{"salir", "salgo", "sal"}},
			{word.Verb, []string{"coger", "coge", "llevar", "recoger", "recoge", "pillar"}},
			{word.Verb, []string{"dejar", "dejo", "soltar", "suelta", "suelto"}},
			{word.Verb, []string{"quitar", "quito", "quita"}},
			{word.Verb, []string{"poner", "pongo", "pone"}},
			{word.Verb, []string{"abrir", "abro", "abre"}},
			{word.Verb, []string{"cerrar", "cierro", "cierra"}},
			{word.Verb, []string{"fin", "terminar", "acabar", "sistema"}},
			{word.Verb, []string{"quill", "thenewquill", "tnq"}},
			{word.Verb, []string{"ad"}},
			{word.Verb, []string{"decir", "di", "hablar", "habla"}},
			{word.Verb, []string{"ir", "voy", "vamos", "ve"}},
			{word.Verb, []string{"ex", "exam", "examinar", "examina", "mirar", "mira", "miro"}},
			{word.Verb, []string{"save", "grabar", "graba", "grabo", "salvar", "salva", "salvo"}},
			{word.Verb, []string{"ram", "ramsave"}},
			{word.Noun, []string{"jugador"}},
			{word.Noun, []string{"enano"}},
			{word.Noun, []string{"elfo"}},
			{word.Noun, []string{"norte", "n", "adelante"}},
			{word.Noun, []string{"sur", "s", "atrás"}},
			{word.Noun, []string{"este", "e"}},
			{word.Noun, []string{"oeste", "o", "w"}},
			{word.Noun, []string{"abrigo", "chaqueta"}},
			{word.Noun, []string{"guantes"}},
			{word.Noun, []string{"cofre", "arcón"}},
			{word.Noun, []string{"monedas", "dinero"}},
			{word.Noun, []string{"denario"}},
			{word.Noun, []string{"carta", "papel"}},
			{word.Noun, []string{"antorcha", "tea", "linterna"}},
			{word.Noun, []string{"llave"}},
			{word.Noun, []string{"cinturón"}},
			{word.Noun, []string{"petaca", "botella"}},
			{word.Noun, []string{"talismán", "amuleto"}},
			{word.Noun, []string{"ropa", "vestido"}},
			{word.Noun, []string{"bolsa", "saco", "petate"}},
			{word.Adjective, []string{"encendida"}},
			{word.Adjective, []string{"apagada"}},
			{word.Adjective, []string{"dorada"}},
			{word.Adjective, []string{"plateada"}},
			{word.Pronoun, []string{"el", "la", "los", "las"}},
			{word.Preposition, []string{"dentro", "adentro"}},
			{word.Conjunction, []string{"enton", "luego", "tras", "y"}},
		}

		for _, tc := range testCases {
			require.GreaterOrEqual(t, len(tc.syns), 1)
			labelName := tc.syns[0]

			t.Run(labelName, func(t *testing.T) {
				w, err := a.Words.FindByLabel(labelName)
				require.NoError(t, err, "error finding word by label %q", labelName)
				require.NotNil(t, w)
				assert.Equal(t, tc.kind, w.Type)
				assert.Equal(t, tc.syns, w.Synonyms, "synonyms for %s doesn't match", labelName)
				for _, syn := range tc.syns {
					wFromSyn, err := a.Words.FirstOfAny(syn)
					require.NoError(t, err, "cannot find word %q from synonym %q", labelName, syn)
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
				require.NoError(t, err, "cannot find character %q", tc.label)
				require.NotNil(t, c)

				assert.Equal(t, tc.desc, c.Description)

				noun, err := a.Words.FindByLabel(tc.noun)
				require.NoError(t, err, "cannot find noun %q", tc.noun)
				assert.Equal(t, c.NounID, noun.ID)
				assert.Equal(t, word.Noun, noun.Type)

				if tc.adj != db.UnderscoreLabel.Name {
					adj, err := a.Words.FindByLabel(tc.adj)
					require.NoError(t, err, "cannot find adjective %q", tc.adj)
					assert.Equal(t, c.AdjectiveID, adj.ID)
					assert.Equal(t, word.Adjective, adj.Type)
				}

				loc, err := a.Locations.FindByLabel(tc.locationLabel)
				require.NoError(t, err)
				assert.Equal(t, loc.ID, c.LocationID)

				assert.Equal(t, tc.created, c.Created)
				assert.Equal(t, tc.human, c.Human)

				for k, v := range tc.vars {
					actual, err := a.Variables.FindByLabel(tc.label, k)
					require.NoError(t, err)

					switch val := v.(type) {
					case int:
						assert.Equal(t, val, actual.Int())
					case float32, float64:
						assert.Equal(t, val, actual.Float())
					case bool:
						assert.Equal(t, val, actual.Bool())
					case string:
						assert.Equal(t, val, actual.String())
					default:
						assert.Equal(t, fmt.Sprint(v), actual.String())
					}
				}
			})
		}
	})
}
