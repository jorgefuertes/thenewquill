package compiler_test

import (
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
	})
}
