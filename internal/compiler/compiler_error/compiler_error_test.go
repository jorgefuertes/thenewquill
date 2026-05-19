package compiler_error_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	cerr "github.com/jorgefuertes/thenewquill/internal/compiler/compiler_error"
	"github.com/jorgefuertes/thenewquill/internal/compiler/line"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildersChain(t *testing.T) {
	e := cerr.ErrInvalidLabel.
		WithSection(kind.Item).
		WithFilename("ao.adv").
		WithStack([]line.Line{line.New("excalibur: sword shiny", 12)}).
		WithLine(line.New("excalibur: sword shiny", 12)).
		AddMsg("nested issue").
		AddMsgf("looking for %q", "foo").
		AddErr(errors.New("wrapped"))

	assert.Contains(t, e.Error(), "invalid label",
		"Error() returns the head message")
	assert.True(t, e.Is(cerr.ErrInvalidLabel),
		"Is matches the originating sentinel")
	assert.True(t, e.Is(errors.New("nested issue")),
		"Is matches an added message too")
	assert.True(t, e.Is(errors.New(`looking for "foo"`)),
		"AddMsgf formats and stores the message")
	assert.True(t, e.Is(errors.New("wrapped")),
		"AddErr stores the wrapped error message")
}

func TestIsOK(t *testing.T) {
	assert.True(t, cerr.OK.IsOK())
	assert.False(t, cerr.ErrInvalidLabel.IsOK())
}

func TestDumpValidationError(t *testing.T) {
	buf := &bytes.Buffer{}
	out := cerr.NewOutput("VALIDATION ERROR")
	out.SetErrOutput(buf)

	// Reach Dump via the public path: build an ErrValidation with several
	// messages and let it render.
	e := cerr.ErrValidation.AddMsg("LabelID is required")
	e = e.AddMsg("NounID is required")

	// Dump writes to its own Output instance internally; redirect by hand
	// using a known-good Output.
	out.SetErrOutput(buf)
	require.NotNil(t, out)

	// Run Dump and capture the rendered text. Dump uses os.Stderr by default,
	// so we just verify it doesn't panic and Error() exposes the first msg.
	assert.NotPanics(t, func() { e.Dump() })
	assert.Equal(t, "validation error", e.Error())
}

func TestDumpRegularErrorDoesNotPanic(t *testing.T) {
	e := cerr.ErrInvalidLabel.
		WithSection(kind.Item).
		WithFilename("a.adv").
		WithStack([]line.Line{
			line.New("a", 1),
			line.New("b", 2),
			line.New("c", 3),
		}).
		WithLine(line.New("c", 3)).
		AddMsg("more detail")

	assert.NotPanics(t, func() { e.Dump() })
}

func TestSetErrOutputRedirectsPrint(t *testing.T) {
	buf := &bytes.Buffer{}

	out := cerr.NewOutput("MY TITLE")
	out.SetErrOutput(buf)
	out.Print()

	rendered := buf.String()
	assert.NotEmpty(t, rendered)
	assert.True(t, strings.Contains(rendered, "MY TITLE"),
		"the title should be rendered inside the bordered output")
}
