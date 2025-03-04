package item_test

import (
	"testing"

	"thenewquill/internal/adventure/item"
	"thenewquill/internal/adventure/words"
	"thenewquill/internal/util"

	"github.com/stretchr/testify/require"
)

func TestItem(t *testing.T) {
	name1 := &words.Word{Label: "name1"}
	adj1 := &words.Word{Label: "adj1"}

	t.Run("new item", func(t *testing.T) {
		i := item.New("test", name1, adj1)
		require.Equal(t, i.Label, "test")
	})

	t.Run("total weight", func(t *testing.T) {
		container := item.New("container", name1, adj1)
		container.IsContainer = true
		container.Weight = 7
		container.MaxWeight = 27
		total := 7
		for range 10 {
			newObject := item.New(util.RandomString(16), name1, adj1)
			newObject.Weight = 2
			require.NoError(t, container.Put(newObject))
			total += 2
		}

		w := container.WeightTotal()
		require.Equal(t, total, w)
	})

	t.Run("over weight", func(t *testing.T) {
		container := item.New("container", name1, adj1)
		container.IsContainer = true
		container.Weight = 5
		container.MaxWeight = 10

		it := item.New("weighted item", name1, adj1)
		it.Weight = 10
		err := container.Put(it)
		require.Error(t, err)
		require.Equal(t, err, item.ErrContainerCantCarrySoMuch)

		it2 := item.New("weighted item 2", name1, adj1)
		it2.Weight = 5
		err = container.Put(it2)
		require.NoError(t, err)

		it3 := item.New("weighted item 3", name1, adj1)
		it3.Weight = 5
		err = container.Put(it3)
		require.Equal(t, err, item.ErrContainerIsFull)
	})

	t.Run("can't put an item in an item which is not a container", func(t *testing.T) {
		it := item.New("my item", name1, adj1)
		it.Weight = 10
		it2 := item.New("another item", name1, adj1)
		it2.Weight = 10
		err := it.Put(it2)
		require.Error(t, err)
		require.Equal(t, err, item.ErrNotContainer)
	})
}
