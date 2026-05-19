package item_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/item"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuery(t *testing.T) {
	svc, _, labels := newTestService(t)

	sword := mustCreate(t, svc, func(i *item.Item) {
		i.LabelID = labels["sword"]
		i.NounID = labels["noun-sword"]
		i.Description = "sword"
	})
	shield := mustCreate(t, svc, func(i *item.Item) {
		i.LabelID = labels["shield"]
		i.NounID = labels["noun-shield"]
		i.Description = "shield"
	})
	bag := mustCreate(t, svc, func(i *item.Item) {
		i.LabelID = labels["bag"]
		i.NounID = labels["noun-bag"]
		i.Description = "bag"
		i.Container = true
	})

	t.Run("Count", func(t *testing.T) {
		assert.Equal(t, 3, svc.Get().Count())
	})

	t.Run("Exists by label", func(t *testing.T) {
		assert.True(t, svc.Get().WithLabel("sword").Exists())
		assert.False(t, svc.Get().WithLabel("nonexistent").Exists())
	})

	t.Run("Exists by label id", func(t *testing.T) {
		assert.True(t, svc.Get().WithLabelID(labels["sword"]).Exists())
	})

	t.Run("First by ID", func(t *testing.T) {
		got, err := svc.Get().WithID(sword.GetID()).First()
		require.NoError(t, err)
		assert.Equal(t, "sword", got.Description)
	})

	t.Run("First missing returns error", func(t *testing.T) {
		_, err := svc.Get().WithID(9999).First()
		require.Error(t, err)
	})

	t.Run("WithNoID excludes", func(t *testing.T) {
		// Three items total; excluding sword leaves shield and bag.
		assert.Equal(t, 2, svc.Get().WithNoID(sword.GetID()).Count())

		got, err := svc.Get().WithNoID(sword.GetID()).WithLabel("shield").First()
		require.NoError(t, err)
		assert.Equal(t, shield.GetID(), got.GetID())
	})

	t.Run("WithLabel returns the right item", func(t *testing.T) {
		got, err := svc.Get().WithLabel("bag").First()
		require.NoError(t, err)
		assert.Equal(t, bag.GetID(), got.GetID())
		assert.True(t, got.Container)
	})
}
