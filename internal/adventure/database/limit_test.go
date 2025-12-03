package database

import (
	"math"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabelLimit(t *testing.T) {
	t.Run("Reach limit", func(t *testing.T) {
		database := New()
		database.nextID = math.MaxUint32

		s1 := &testStoreable{ID: primitive.UndefinedID, Title: "test-0"}
		_, err := database.Create(s1)
		require.NoError(t, err)

		s2 := &testStoreable{ID: primitive.UndefinedID, Title: "too-many-records"}
		_, err = database.Create(s2)
		assert.ErrorIs(t, err, ErrLimitReached)
	})
}
