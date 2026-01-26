package compiler_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/compiler"

	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"

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

	t.Run("wrong", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := compiler.Compile(tc.filename)
				require.Error(t, err)
				er, ok := err.(cerr.CompilerError)
				require.True(t, ok, "error is not a CompilerError")
				require.True(
					t,
					er.Is(tc.expectedError),
					"error is not the expected %q, actiual: %q",
					tc.expectedError,
					er.Error(),
				)
			})
		}
	})
}
