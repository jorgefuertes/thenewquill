package compiler_test

import (
	"io"
	"testing"

	"github.com/fatih/color"
	"github.com/jorgefuertes/thenewquill/internal/compiler"
	"github.com/jorgefuertes/thenewquill/pkg/log"

	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"

	"github.com/stretchr/testify/require"
)

func TestWrongFiles(t *testing.T) {
	whiteOnRed := color.New(color.FgWhite, color.BgRed).SprintFunc()
	t.Log(whiteOnRed("This text should be white on red"))

	testCases := []struct {
		name          string
		filename      string
		expectedError error
	}{
		{"unclosed comment", "src/wrong/unclosed_comment.adv", cerr.ErrUnclosedComment},
		{"unclosed string", "src/wrong/unclosed_string.adv", cerr.ErrUnclosedMultiline},
		{"duplicated synonyms", "src/wrong/duped_syn.adv", cerr.ErrDuplicatedSynonym},
	}

	// disable logging
	log.SetOutput(io.Discard)

	t.Run("wrong", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := compiler.Compile(tc.filename)
				require.Error(t, err)
				er, ok := err.(cerr.CompilerError)
				require.True(t, ok, "error is not a CompilerError")
				require.True(t, er.Is(tc.expectedError), "error is not the expected one: %s", er.Error())
			})
		}
	})
}
