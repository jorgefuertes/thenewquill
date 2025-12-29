package variable_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/variable"
	"github.com/jorgefuertes/thenewquill/internal/util"
	"github.com/stretchr/testify/assert"
)

const nonConvertible = "non-convertible"

func TestVariable(t *testing.T) {
	t.Run("SetID", func(t *testing.T) {
		v := &variable.Variable{Value: "2"}
		v.SetID(5)
		assert.EqualValues(t, 5, v.ID)
		v.SetID(6)
		assert.EqualValues(t, 6, v.ID)
	})

	t.Run("SetLabelID", func(t *testing.T) {
		v := &variable.Variable{Value: "2"}
		v.SetLabelID(5)
		assert.EqualValues(t, 5, v.LabelID)
		v.SetLabelID(6)
		assert.EqualValues(t, 6, v.LabelID)
	})

	t.Run("int", func(t *testing.T) {
		v := &variable.Variable{Value: "2"}
		assert.Equal(t, kind.Variable, kind.KindOf(v))
		assert.Equal(t, 2, v.Int())
		assert.Equal(t, "2", v.String())

		v.SetValue(byte(6))
		assert.Equal(t, 6, v.Int())

		v.SetValue("8")
		assert.Equal(t, 8, v.Int())

		v.SetValue(9.1)
		assert.Equal(t, 9.1, v.Float())

		v.SetValue([]byte(nonConvertible))
		assert.Equal(t, 0, v.Int())

		v.SetValue(true)
		assert.Equal(t, 1, v.Int())
		assert.True(t, v.Bool())
		assert.True(t, v.IsTrue())
		assert.False(t, v.IsFalse())

		v.SetValue(false)
		assert.Equal(t, 0, v.Int())
		assert.False(t, v.Bool())
		assert.True(t, v.IsFalse())
		assert.False(t, v.IsTrue())
	})

	t.Run("Float", func(t *testing.T) {
		v := &variable.Variable{Value: "2.5"}
		assert.Equal(t, 2.5, v.Float())

		v.SetValue(true)
		assert.Equal(t, 1.0, v.Float())

		v.SetValue(false)
		assert.Equal(t, 0.0, v.Float())

		v.SetValue("2.5")
		assert.Equal(t, 2.5, v.Float())

		v.SetValue([]byte(nonConvertible))
		assert.Equal(t, 0.0, v.Float())
	})

	t.Run("String", func(t *testing.T) {
		v := &variable.Variable{Value: "test-string"}
		assert.Equal(t, "test-string", v.String())

		v.SetValue(2)
		assert.Equal(t, "2", v.String())

		v.SetValue(true)
		assert.Equal(t, "1", v.String())

		v.SetValue(false)
		assert.Equal(t, "0", v.String())

		v.SetValue(2.5445)
		assert.Equal(t, "2.5445", v.String())

		v.SetValue(2.54459999999)
		assert.Equal(t, "2.5446", v.String())

		v.SetValue(2.544511111111)
		assert.Equal(t, "2.5445", v.String())

		v.SetValue(byte('a'))
		assert.Equal(t, "97", v.String())

		v.SetValue([]byte("test bytes"))
		assert.Equal(t, "[116 101 115 116 32 98 121 116 101 115]", v.String())

		v.SetValue([]int{1, 2, 3})
		assert.Equal(t, "[1 2 3]", v.String())
	})

	t.Run("Bool", func(t *testing.T) {
		v := &variable.Variable{Value: util.ValueToString(true)}
		assert.Equal(t, true, v.Bool())

		v.SetValue(2)
		assert.Equal(t, true, v.Bool())

		v.SetValue(true)
		assert.Equal(t, true, v.Bool())

		v.SetValue(false)
		assert.Equal(t, false, v.Bool())

		v.SetValue(2.5445)
		assert.Equal(t, true, v.Bool())

		v.SetValue(byte('a'))
		assert.Equal(t, true, v.Bool())

		v.SetValue([]byte(nonConvertible))
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
				v.SetValue(tc.val)
				assert.Equal(t, tc.want, v.Bool(), "val: %s", tc.val)
			})
		}
	})
}
