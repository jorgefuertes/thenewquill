package item_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/item"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsContained(t *testing.T) {
	svc, _, labels := newTestService(t)

	bag := mustCreate(t, svc, func(i *item.Item) {
		i.LabelID = labels["bag"]
		i.NounID = labels["noun-bag"]
		i.Description = "bag"
		i.Container = true
	})

	t.Run("At==0 means not contained", func(t *testing.T) {
		i := *item.New()
		assert.False(t, svc.IsContained(i))
	})

	t.Run("At points to a non-item id means not contained", func(t *testing.T) {
		i := *item.New()
		i.At = 9999 // no item has this id
		assert.False(t, svc.IsContained(i))
	})

	t.Run("At points to a container item means contained", func(t *testing.T) {
		coin := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["coin"]
			i.NounID = labels["noun-coin"]
			i.Description = "coin"
			i.At = bag.GetID()
		})

		assert.True(t, svc.IsContained(*coin))
	})
}

func TestGetItemContainer(t *testing.T) {
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

	t.Run("returns the container item", func(t *testing.T) {
		got, err := svc.GetItemContainer(*coin)
		require.NoError(t, err)
		assert.Equal(t, bag.GetID(), got.GetID())
	})

	t.Run("returns error for unknown container id", func(t *testing.T) {
		orphan := item.New()
		orphan.At = 99999

		_, err := svc.GetItemContainer(*orphan)
		require.Error(t, err)
	})
}

func TestPutInto(t *testing.T) {
	t.Run("puts an item into a container", func(t *testing.T) {
		svc, _, labels := newTestService(t)

		bag := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["bag"]
			i.NounID = labels["noun-bag"]
			i.Description = "bag"
			i.Container = true
			i.MaxWeight = 50
		})
		coin := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["coin"]
			i.NounID = labels["noun-coin"]
			i.Description = "coin"
			i.Weight = 1
		})

		require.NoError(t, svc.PutInto(coin, *bag))

		got, err := svc.Get().WithID(coin.GetID()).First()
		require.NoError(t, err)
		assert.Equal(t, bag.GetID(), got.At)
	})

	t.Run("rejects an already-contained item", func(t *testing.T) {
		svc, _, labels := newTestService(t)

		bag1 := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["big-bag"]
			i.NounID = labels["noun-bag"]
			i.Description = "bag1"
			i.Container = true
			i.MaxWeight = 50
		})
		bag2 := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["small-bag"]
			i.NounID = labels["noun-chest"]
			i.Description = "bag2"
			i.Container = true
			i.MaxWeight = 50
		})
		coin := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["coin"]
			i.NounID = labels["noun-coin"]
			i.Description = "coin"
			i.At = bag1.GetID()
		})

		err := svc.PutInto(coin, *bag2)
		require.ErrorIs(t, err, item.ErrItemAlreadyContained)
	})

	t.Run("rejects when container would exceed its max weight", func(t *testing.T) {
		svc, _, labels := newTestService(t)

		bag := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["bag"]
			i.NounID = labels["noun-bag"]
			i.Description = "small bag"
			i.Container = true
			i.MaxWeight = 1
		})
		stone := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["sword"]
			i.NounID = labels["noun-sword"]
			i.Description = "heavy stone"
			i.Weight = 100
		})

		err := svc.PutInto(stone, *bag)
		require.ErrorIs(t, err, item.ErrContainerCantCarrySoMuch)
	})
}

func TestContents(t *testing.T) {
	svc, _, labels := newTestService(t)

	bag := mustCreate(t, svc, func(i *item.Item) {
		i.LabelID = labels["bag"]
		i.NounID = labels["noun-bag"]
		i.Description = "bag"
		i.Container = true
	})
	mustCreate(t, svc, func(i *item.Item) {
		i.LabelID = labels["coin"]
		i.NounID = labels["noun-coin"]
		i.Description = "coin"
		i.At = bag.GetID()
	})
	mustCreate(t, svc, func(i *item.Item) {
		i.LabelID = labels["ring"]
		i.NounID = labels["noun-ring"]
		i.Description = "ring"
		i.At = bag.GetID()
	})
	// An item outside the bag should not appear.
	mustCreate(t, svc, func(i *item.Item) {
		i.LabelID = labels["sword"]
		i.NounID = labels["noun-sword"]
		i.Description = "sword"
		i.At = 99
	})

	contents := svc.Contents(bag.GetID())
	assert.Len(t, contents, 2)
}

func TestTotalWeight(t *testing.T) {
	svc, _, labels := newTestService(t)

	t.Run("non-container returns its own weight", func(t *testing.T) {
		i := item.New()
		i.Weight = 7

		assert.Equal(t, 7, svc.TotalWeight(*i))
	})

	t.Run("container with only leaf items sums shell + leaves", func(t *testing.T) {
		bag := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["bag"]
			i.NounID = labels["noun-bag"]
			i.Description = "bag"
			i.Container = true
			i.Weight = 1
		})
		mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["coin"]
			i.NounID = labels["noun-coin"]
			i.Description = "coin"
			i.Weight = 2
			i.At = bag.GetID()
		})
		mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["ring"]
			i.NounID = labels["noun-ring"]
			i.Description = "ring"
			i.Weight = 3
			i.At = bag.GetID()
		})

		assert.Equal(t, 6, svc.TotalWeight(*bag))
	})

	t.Run("container sums shell + nested mixed contents", func(t *testing.T) {
		outer := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["big-bag"]
			i.NounID = labels["noun-bag"]
			i.Description = "outer"
			i.Container = true
			i.Weight = 2
		})
		inner := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["small-bag"]
			i.NounID = labels["noun-chest"]
			i.Description = "inner"
			i.Container = true
			i.Weight = 5
			i.At = outer.GetID()
		})
		mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["coin"]
			i.NounID = labels["noun-coin"]
			i.Description = "coin inside outer"
			i.Weight = 1
			i.At = outer.GetID()
		})
		mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["ring"]
			i.NounID = labels["noun-ring"]
			i.Description = "ring inside inner"
			i.Weight = 4
			i.At = inner.GetID()
		})

		// outer (2) + coin (1) + inner (5) + ring (4) = 12
		assert.Equal(t, 12, svc.TotalWeight(*outer))
	})
}
