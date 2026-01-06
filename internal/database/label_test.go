package database_test

import (
	"strings"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseLabels(t *testing.T) {
	db := database.NewDB()

	testCases := []struct {
		labelName         string
		expectedID        uint32
		expectedCreateErr error
		expectedGetErr    error
	}{
		{"dewarf", 3, nil, nil},
		{"key", 4, nil, nil},
		{"light-bulb", 5, nil, nil},
		{"lightbulb", 6, nil, nil},
		{"LightBulb", 6, nil, nil},
		{"dewarf", 3, nil, nil},
		{"light.on", 7, nil, nil},
		{"light.off", 8, nil, nil},
		{"light-on", 9, nil, nil},
		{"#undefined", 0, database.ErrInvalidLabel, nil},
		{"test!!!", 0, database.ErrInvalidLabel, database.ErrLabelNotFound},
		{"*", 1, nil, nil},
		{"_", 2, nil, nil},
	}

	for _, tc := range testCases {
		t.Run("create/"+tc.labelName, func(t *testing.T) {
			id, err := db.CreateLabel(tc.labelName)
			assert.Equal(t, tc.expectedCreateErr, err)
			assert.Equal(t, tc.expectedID, id)
		})
	}

	for _, tc := range testCases {
		t.Run("get/"+tc.labelName, func(t *testing.T) {
			id, err := db.GetLabelID(tc.labelName)
			assert.Equal(t, tc.expectedGetErr, err)
			assert.Equal(t, tc.expectedID, id)

			if tc.expectedGetErr != nil {
				return
			}

			label, err := db.GetLabel(id)
			assert.Equal(t, tc.expectedGetErr, err)
			assert.Equal(t, strings.ToLower(tc.labelName), label)
		})
	}
}
