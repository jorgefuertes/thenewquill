package compiler_test

import (
	"io"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/compiler"
	log "github.com/jorgefuertes/thenewquill/internal/log"

	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"

	"github.com/stretchr/testify/require"
)

func TestWrongFiles(t *testing.T) {
	testCases := []struct {
		name           string
		filename       string
		expectedErrors []error
	}{
		{"unclosed comment", "src/wrong/unclosed_comment.adv", []error{cerr.ErrUnclosedComment}},
		{"unclosed string", "src/wrong/unclosed_string.adv", []error{cerr.ErrUnclosedMultiline}},
		{"duplicated synonyms", "src/wrong/duped_syn.adv", []error{cerr.ErrDuplicatedSynonym, character.ErrNoHuman}},
	}

	// disable logging
	log.SetOutput(io.Discard)

	t.Run("test wrong files", func(t *testing.T) {
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := compiler.Compile(tc.filename)
				require.Error(t, err)
				for _, er := range tc.expectedErrors {
					require.ErrorContains(t, err, er.Error())
				}
			})
		}
	})
}
