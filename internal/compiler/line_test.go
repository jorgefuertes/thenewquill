package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLine(t *testing.T) {
	l := line{text: `test:    "This is  a test \"message\"."`, n: 1}
	assert.Equal(t, `test:    "This is  a test \"message\"."`, l.optimized())

	l = line{text: `test: "This is a test \"message\"." // with a comment`, n: 1}
	assert.Equal(t, `test: "This is a test \"message\"."`, l.optimized())

	l = line{text: `test: "This is a test \"message\"." /* with a comment */`, n: 1}
	assert.Equal(t, `test: "This is a test \"message\"."`, l.optimized())

	l = line{text: `test: "This is a test message`, n: 1}
	_, _, ok := l.toVar()
	assert.False(t, ok)

	l = line{text: `test = "This is a test message"`, n: 1}
	name, value, ok := l.toVar()
	assert.True(t, ok)
	assert.Equal(t, "test", name)
	assert.Equal(t, "This is a test message", value)

	// TODO: check for balanced quotes
	l = line{text: `test = "This is a test message`, n: 1}
	name, value, ok = l.toVar()
	assert.True(t, ok)
	assert.Equal(t, "test", name)
	assert.Equal(t, "This is a test message", value)
}
