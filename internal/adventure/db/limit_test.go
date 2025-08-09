package db

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabelLimit(t *testing.T) {
	t.Run("Reach limit", func(t *testing.T) {
		database := New()

		_, err := database.AddLabel("test-0", false)
		require.NoError(t, err)

		database.nextID = math.MaxUint16

		_, err = database.AddLabel("label-too-many", false)
		assert.ErrorIs(t, err, ErrLimitReached)
	})
}
