package database

import (
	"fmt"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database/adapter"
	"github.com/stretchr/testify/require"
)

type testValue struct {
	ID      uint32
	LabelID uint32
	Value   string
}

var _ adapter.Storeable = &testValue{}

func (t testValue) GetKind() kind.Kind {
	return kind.Test
}

func (t testValue) GetID() uint32 {
	return t.ID
}

func (t *testValue) SetID(id uint32) {
	t.ID = id
}

func (t *testValue) SetLabelID(id uint32) {
	t.LabelID = id
}

func (t *testValue) GetLabelID() uint32 {
	return t.LabelID
}

func TestCursor(t *testing.T) {
	const limit = 10

	c := newCursor()

	t.Run("addOrReplace", func(t *testing.T) {
		for i := 1; i <= limit; i++ {
			td := testValue{
				ID:      uint32(i),
				LabelID: uint32(i + 2),
				Value:   fmt.Sprintf("test-data-%d", i),
			}

			b, err := cbor.Marshal(td)
			require.NoError(t, err)

			c.addOrReplace(Record{LabelID: td.LabelID, Kind: kind.Test, Data: b})
		}

		t.Run("Count", func(t *testing.T) {
			require.Equal(t, limit, c.Count())
		})

		t.Run("Exists", func(t *testing.T) {
			require.True(t, c.Exists())
		})

		t.Run("First", func(t *testing.T) {
			var td testValue
			require.NoError(t, c.First(&td))
			require.NotZero(t, td.ID)
			require.NotZero(t, td.LabelID)
			require.Equal(t, fmt.Sprintf("test-data-%d", td.ID), td.Value)
		})

		t.Run("getByIndex", func(t *testing.T) {
			for i := 0; i < c.Count(); i++ {
				r, ok := c.getByIndex(i)
				require.True(t, ok)
				require.NotZero(t, r.LabelID)
				require.NotEmpty(t, r.Data)
			}
		})

		t.Run("Next", func(t *testing.T) {
			items := make(map[uint32]testValue, 0)

			for i := 1; i <= limit; i++ {
				var td testValue

				require.True(t, c.Next(&td))
				t.Logf("item %d index %d", td.ID, c.i)

				require.NotZero(t, td.ID)
				require.Equal(t, td.ID+2, td.LabelID)
				require.Equal(t, fmt.Sprintf("test-data-%d", td.ID), td.Value)

				_, ok := items[td.ID]
				require.False(t, ok, "item %d already exists", td.ID)

				items[td.ID] = td
			}

			require.Len(t, items, limit)
		})
	})
}
