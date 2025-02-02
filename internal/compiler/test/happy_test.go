package compiler_test

import (
	"testing"

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
}
