package item_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/item"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsAt(t *testing.T) {
	svc := &item.Service{}

	i := item.New()
	i.At = 42

	assert.True(t, svc.IsAt(*i, 42))
	assert.False(t, svc.IsAt(*i, 99))
	assert.False(t, svc.IsAt(*i, 0))
}

func TestMoveTo(t *testing.T) {
	t.Run("moves an unowned item", func(t *testing.T) {
		svc, _, labels := newTestService(t)

		sword := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["sword"]
			i.NounID = labels["noun-sword"]
			i.Description = "sword"
		})

		require.NoError(t, svc.MoveTo(sword, 17))

		got, err := svc.Get().WithID(sword.GetID()).First()
		require.NoError(t, err)
		assert.Equal(t, uint32(17), got.At)
	})

	t.Run("re-moves an item already at a location", func(t *testing.T) {
		// An item placed at a location (not inside a container) should be
		// movable to another location.
		svc, _, labels := newTestService(t)

		sword := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["sword"]
			i.NounID = labels["noun-sword"]
			i.Description = "sword"
			i.At = 100 // a location id, not a container item id
		})

		require.NoError(t, svc.MoveTo(sword, 200))

		got, err := svc.Get().WithID(sword.GetID()).First()
		require.NoError(t, err)
		assert.Equal(t, uint32(200), got.At)
	})

	t.Run("blocked when item is inside a container", func(t *testing.T) {
		svc, _, labels := newTestService(t)

		bag := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["bag"]
			i.NounID = labels["noun-bag"]
			i.Description = "bag"
			i.Container = true
		})

		coin := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["coin"]
			i.NounID = labels["noun-coin"]
			i.Description = "coin"
			i.At = bag.GetID()
		})

		err := svc.MoveTo(coin, 999)
		require.ErrorIs(t, err, item.ErrItemAlreadyContained)
	})
}
