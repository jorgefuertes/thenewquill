package variable_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/variable"
	"github.com/stretchr/testify/assert"
)

func TestVariable(t *testing.T) {
	t.Run("SetID", func(t *testing.T) {
		v := variable.New(2)
		s := v.SetID(db.ID(5))
		assert.Equal(t, db.ID(5), s.GetID())
	})

	t.Run("int", func(t *testing.T) {
		v := variable.New(2)
		assert.Equal(t, db.Variables, v.GetKind())
		assert.Equal(t, 2, v.Value)
		assert.Equal(t, 2, v.Int())
		assert.Equal(t, 2.0, v.Float())

		v.Value = byte(6)
		assert.Equal(t, 6, v.Int())

		v.Value = "8"
		assert.Equal(t, 8, v.Int())

		v.Value = 9.1
		assert.Equal(t, 9, v.Int())

		v.Value = []byte("non-convertible")
		assert.Equal(t, 0, v.Int())

		v.Value = true
		assert.Equal(t, 1, v.Int())

		v.Value = false
		assert.Equal(t, 0, v.Int())
	})

	t.Run("Float", func(t *testing.T) {
		v := variable.New(2.5)
		assert.Equal(t, 2.5, v.Float())

		v.Value = true
		assert.Equal(t, 1.0, v.Float())

		v.Value = false
		assert.Equal(t, 0.0, v.Float())

		v.Value = "2.5"
		assert.Equal(t, 2.5, v.Float())

		v.Value = []byte("non-convertible")
		assert.Equal(t, 0.0, v.Float())
	})

	t.Run("String", func(t *testing.T) {
		v := variable.New("test-string")
		assert.Equal(t, "test-string", v.String())

		v.Value = 2
		assert.Equal(t, "2", v.String())

		v.Value = true
		assert.Equal(t, "true", v.String())

		v.Value = false
		assert.Equal(t, "false", v.String())

		v.Value = 2.5445
		assert.Equal(t, "2.54", v.String())

		v.Value = byte('a')
		assert.Equal(t, "a", v.String())

		v.Value = []byte("test bytes")
		assert.Equal(t, "test bytes", v.String())

		v.Value = []int{1, 2, 3}
		assert.Equal(t, "[1 2 3]", v.String())
	})

	t.Run("Bool", func(t *testing.T) {
		v := variable.New(true)
		assert.Equal(t, true, v.Bool())

		v.Value = 2
		assert.Equal(t, true, v.Bool())

		v.Value = true
		assert.Equal(t, true, v.Bool())

		v.Value = false
		assert.Equal(t, false, v.Bool())

		v.Value = 2.5445
		assert.Equal(t, true, v.Bool())

		v.Value = byte('a')
		assert.Equal(t, true, v.Bool())

		v.Value = []byte("non-convertible")
		assert.Equal(t, false, v.Bool())

		testCases := []struct {
			val  string
			want bool
		}{
			{"s", true},
			{"Si", true},
			{"sí", true},
			{"yes", true},
			{"YES", true},
			{"true", true},
			{"TRUE", true},
			{"on", true},
			{"ON", true},
			{"no", false},
			{"NO", false},
			{"false", false},
			{"FALSE", false},
			{"off", false},
			{"OFF", false},
			{"-1", false},
			{"0", false},
			{"1", true},
			{"2", false},
			{"t", true},
			{"f", false},
			{"T", true},
			{"F", false},
		}

		for _, tc := range testCases {
			v.Value = tc.val
			assert.Equal(t, tc.want, v.Bool(), "val: %s", tc.val)
		}
	})

	t.Run("Byte", func(t *testing.T) {
		testCases := []struct {
			val  any
			want byte
		}{
			{32, 32},
			{" ", 32},
			{true, 1},
			{false, 0},
			{"true", 1},
			{"false", 0},
			{"abcdef", 97},
			{" ", 32},
			{[]int{1, 2, 3}, 0},
		}

		for _, tc := range testCases {
			v := variable.New(tc.val)
			assert.Equal(t, tc.want, v.Byte(), "val: %v", tc.val)
		}
	})
}
