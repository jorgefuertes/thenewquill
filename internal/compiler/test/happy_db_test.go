package compiler_test

import (
	"bytes"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/compiler"

	"github.com/stretchr/testify/require"
)

func TestCompilerDB(t *testing.T) {
	const srcFilename = "src/happy/test.adv"

	var a1 *adventure.Adventure

	t.Run("compile", func(t *testing.T) {
		var err error

		a1, err = compiler.Compile(srcFilename)
		require.NoError(t, err)
		require.NotNil(t, a1)

		t.Run("adventure to buffer", func(t *testing.T) {
			saveBuf := bytes.NewBuffer(nil)
			err = a1.Export(saveBuf)
			require.NoError(t, err)
			require.NotNil(t, saveBuf)
			require.NotZero(t, saveBuf.Len())

			// t.Run("buffer to DB", func(t *testing.T) {
			// 	loadBuf := bytes.NewBuffer(saveBuf.Bytes())
			// 	a2 = adventure.New()
			// 	err := a2.DB.Import(loadBuf)
			// 	require.NoError(t, err)

			// 	h1, err := a1.DB.Hash()
			// 	require.NoError(t, err)
			// 	h2, err := a2.DB.Hash()
			// 	require.NoError(t, err)
			// 	require.Equal(t, h1, h2)
			// })
		})
	})
}
