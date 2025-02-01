package obj_test

import (
	"testing"

	"thenewquill/internal/adventure/obj"
	"thenewquill/internal/adventure/voc"
	"thenewquill/internal/util"

	"github.com/stretchr/testify/require"
)

func TestItem(t *testing.T) {
	name1 := &voc.Word{Label: "name1"}
	adj1 := &voc.Word{Label: "adj1"}

	t.Run("new item", func(t *testing.T) {
		i := obj.New("test", name1, adj1)
		require.Equal(t, i.Label(), "test")
	})

	t.Run("total weight", func(t *testing.T) {
		container := obj.New("container", name1, adj1)
		container.SetContainer()
		container.SetWeight(7)
		container.SetMaxWeight(27)
		total := 7
		for range 10 {
			newObject := obj.New(util.RandomString(16), name1, adj1)
			newObject.SetWeight(2)
			container.Put(newObject)
			total += 2
		}

		w := container.WeightTotal()
		require.Equal(t, total, w)
	})

	t.Run("over weight", func(t *testing.T) {
		container := obj.New("container", name1, adj1)
		container.SetContainer()
		container.SetWeight(5)
		container.SetMaxWeight(10)

		it := obj.New("weighted item", name1, adj1)
		it.SetWeight(10)
		err := container.Put(it)
		require.Error(t, err)
		require.Equal(t, err, obj.ErrContainerCantCarrySoMuch)

		it2 := obj.New("weighted item 2", name1, adj1)
		it2.SetWeight(5)
		err = container.Put(it2)
		require.NoError(t, err)

		it3 := obj.New("weighted item 3", name1, adj1)
		it3.SetWeight(5)
		err = container.Put(it3)
		require.Equal(t, err, obj.ErrContainerIsFull)
	})

	t.Run("can't put an item in an item which is not a container", func(t *testing.T) {
		it := obj.New("my item", name1, adj1)
		it.SetWeight(10)
		it2 := obj.New("another item", name1, adj1)
		it2.SetWeight(10)
		err := it.Put(it2)
		require.Error(t, err)
		require.Equal(t, err, obj.ErrNotContainer)
	})
}
