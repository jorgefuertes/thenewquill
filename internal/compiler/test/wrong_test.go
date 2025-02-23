package compiler_test

import (
	"testing"

	"thenewquill/internal/compiler"

	"github.com/stretchr/testify/require"
)

func TestWrongFiles(t *testing.T) {
	testCases := []struct {
		name          string
		filename      string
		expectedError error
	}{
		{"unclosed comment", "src/wrong/unclosed_comment.adv", compiler.ErrUnclosedComment},
		{"unclosed string", "src/wrong/unclosed_string.adv", compiler.ErrUnclosedMultiline},
	}

	t.Run("test wrong files", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := compiler.Compile(tc.filename)
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedError)
			})
		}
	})
}
