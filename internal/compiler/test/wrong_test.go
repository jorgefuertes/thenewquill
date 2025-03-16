package compiler_test

import (
	"testing"

	"thenewquill/internal/compiler"

	cerr "thenewquill/internal/compiler/compiler_error"

	"github.com/stretchr/testify/require"
)

func TestWrongFiles(t *testing.T) {
	testCases := []struct {
		name          string
		filename      string
		expectedError error
	}{
		{"unclosed comment", "src/wrong/unclosed_comment.adv", cerr.ErrUnclosedComment},
		{"unclosed string", "src/wrong/unclosed_string.adv", cerr.ErrUnclosedMultiline},
		{"duplicated synonyms", "src/wrong/duped_syn.adv", cerr.ErrDuplicatedSynonym},
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
