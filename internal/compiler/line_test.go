package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLine(t *testing.T) {
	l := line{text: `test:    "This is  a test \"message\"."`, n: 1}
	assert.Equal(t, `test: "This is a test \"message\"."`, l.optimized())

	l = line{text: `test: "This is a test \"message\"." // with a comment`, n: 1}
	assert.Equal(t, `test: "This is a test \"message\"."`, l.optimized())

	l = line{text: `test: "This is a test \"message\"." /* with a comment */`, n: 1}
	assert.Equal(t, `test: "This is a test \"message\"."`, l.optimized())
}
