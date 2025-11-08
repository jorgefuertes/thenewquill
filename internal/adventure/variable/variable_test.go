package variable_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/variable"
	"github.com/stretchr/testify/assert"
)

func TestVariable(t *testing.T) {
	t.Run("SetID", func(t *testing.T) {
		v := variable.New(db.ID(5), 2)
		assert.Equal(t, db.ID(5), v.GetID())
		s := v.SetID(6)
		assert.Equal(t, db.ID(6), s.GetID())
	})

	t.Run("int", func(t *testing.T) {
		v := variable.New(0, 2)
		assert.Equal(t, kind.Variable, kind.KindOf(v))
		assert.Equal(t, 2, v.Int())
		assert.Equal(t, "2", v.String())

		v.Set(byte(6))
		assert.Equal(t, 6, v.Int())

		v.Set("8")
		assert.Equal(t, 8, v.Int())

		v.Set(9.1)
		assert.Equal(t, 9.1, v.Float())

		v.Set([]byte("non-convertible"))
		assert.Equal(t, 0, v.Int())

		v.Set(true)
		assert.Equal(t, 1, v.Int())
		assert.True(t, v.Bool())
		assert.True(t, v.IsTrue())
		assert.False(t, v.IsFalse())

		v.Set(false)
		assert.Equal(t, 0, v.Int())
		assert.False(t, v.Bool())
		assert.True(t, v.IsFalse())
		assert.False(t, v.IsTrue())
	})

	t.Run("Float", func(t *testing.T) {
		v := variable.New(1, 2.5)
		assert.Equal(t, 2.5, v.Float())

		v.Set(true)
		assert.Equal(t, 1.0, v.Float())

		v.Set(false)
		assert.Equal(t, 0.0, v.Float())

		v.Set("2.5")
		assert.Equal(t, 2.5, v.Float())

		v.Set([]byte("non-convertible"))
		assert.Equal(t, 0.0, v.Float())
	})

	t.Run("String", func(t *testing.T) {
		v := variable.New(1, "test-string")
		assert.Equal(t, "test-string", v.String())

		v.Set(2)
		assert.Equal(t, "2", v.String())

		v.Set(true)
		assert.Equal(t, "1", v.String())

		v.Set(false)
		assert.Equal(t, "0", v.String())

		v.Set(2.5445)
		assert.Equal(t, "2.5445", v.String())

		v.Set(2.54459999999)
		assert.Equal(t, "2.5446", v.String())

		v.Set(2.544511111111)
		assert.Equal(t, "2.5445", v.String())

		v.Set(byte('a'))
		assert.Equal(t, "97", v.String())

		v.Set([]byte("test bytes"))
		assert.Equal(t, "[116 101 115 116 32 98 121 116 101 115]", v.String())

		v.Set([]int{1, 2, 3})
		assert.Equal(t, "[1 2 3]", v.String())
	})

	t.Run("Bool", func(t *testing.T) {
		v := variable.New(1, true)
		assert.Equal(t, true, v.Bool())

		v.Set(2)
		assert.Equal(t, true, v.Bool())

		v.Set(true)
		assert.Equal(t, true, v.Bool())

		v.Set(false)
		assert.Equal(t, false, v.Bool())

		v.Set(2.5445)
		assert.Equal(t, true, v.Bool())

		v.Set(byte('a'))
		assert.Equal(t, true, v.Bool())

		v.Set([]byte("non-convertible"))
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
			{"2", true},
			{"t", true},
			{"f", false},
			{"T", true},
			{"F", false},
		}

		for _, tc := range testCases {
			t.Run(tc.val, func(t *testing.T) {
				v.Set(tc.val)
				assert.Equal(t, tc.want, v.Bool(), "val: %s", tc.val)
			})
		}
	})
}
