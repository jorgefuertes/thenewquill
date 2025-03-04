package bin_test

import (
	"bytes"
	"testing"

	"thenewquill/internal/adventure"
	"thenewquill/internal/compiler"
	"thenewquill/internal/compiler/bin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeAndDecode(t *testing.T) {
	var a *adventure.Adventure
	var a2 *adventure.Adventure
	w := bytes.NewBuffer(nil)

	t.Run("I can compile the test adventure", func(t *testing.T) {
		var err error
		a, err = compiler.Compile("../test/src/ao/ao.adv")
		require.NoError(t, err)
	})

	t.Run("I can dump the adventure to a buffer", func(t *testing.T) {
		err := bin.Export(a, w)
		require.NoError(t, err)
		require.NotNil(t, w)
		require.NotEmpty(t, w)
	})

	t.Run("I can load the adventure from the buffer", func(t *testing.T) {
		var err error

		r := bytes.NewReader(w.Bytes())
		a2, err = bin.Import(r)
		require.NoError(t, err)
	})

	t.Run("I can compare the adventures", func(t *testing.T) {
		require.Equal(t, a.Config, a2.Config)
		require.Equal(t, a.Vars, a2.Vars)
		require.Equal(t, a.Messages, a2.Messages)
		require.Equal(t, a.Words, a2.Words)
		for _, w := range a.Words {
			w2 := a2.Words.Get(w.Type, w.Label)
			require.NotNil(t, w2)
			assert.True(t, w.IsExactlyEqual(*w2))
		}

		require.Len(t, a.Locations, len(a2.Locations))
		for _, l := range a.Locations {
			l2 := a2.Locations.Get(l.Label)
			require.NotNil(t, l2)
			assert.Equal(t, l.Label, l2.Label)
			assert.Equal(t, l.Title, l2.Title)
			assert.Equal(t, l.Description, l2.Description)
			assert.Len(t, l.Conns, len(l2.Conns))
			for _, c := range l.Conns {
				c2To := l2.GetConn(c.Word)
				require.NotNil(t, c2To, "conn %s not in %+v", c.Word.Label, l2.Conns)
				assert.Equal(t, c.To.Label, c2To.Label)
			}
		}

		require.Len(t, a.Chars, len(a2.Chars))
		for _, c := range a.Chars {
			c2 := a2.Chars.Get(c.Label)
			require.NotNil(t, c2)
			assert.Equal(t, c.Label, c2.Label)
			assert.Equal(t, c.Name.Label, c2.Name.Label)
			assert.Equal(t, c.Adjective.Label, c2.Adjective.Label)
			assert.Equal(t, c.Description, c2.Description)
			assert.Equal(t, c.Location.Label, c2.Location.Label)
			assert.Equal(t, c.Created, c2.Created)
			assert.Equal(t, c.Human, c2.Human)
		}
	})
}
