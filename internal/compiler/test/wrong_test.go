package compiler_test

import (
	"testing"

	"thenewquill/internal/compiler"

	"github.com/stretchr/testify/require"
)

func TestWrongFiles(t *testing.T) {
	filename := "adv_files/wrong/unclosed_comment.adv"

	_, err := compiler.Compile(filename)
	require.Error(t, err)
	require.ErrorIs(t, err, compiler.ErrUnclosedComment)
}
