package util_test

import (
	"testing"

	"thenewquill/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	last := ""
	for range 100 {
		s := util.RandomString(16)
		assert.NotEmpty(t, s)
		assert.Len(t, s, 16)
		assert.NotEqual(t, last, s)
		last = s
	}
}
