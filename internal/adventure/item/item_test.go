package item_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/item"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestService spins up an in-memory database, registers the labels used
// across tests and returns a ready-to-use item service.
func newTestService(t *testing.T) (*item.Service, *database.DB, map[string]uint32) {
	t.Helper()

	db := database.NewDB()
	svc := item.NewService(db)

	labels := map[string]uint32{}
	for _, l := range []string{
		"sword", "shield", "bag", "chest", "coin",
		"big-bag", "small-bag", "ring",
		"noun-sword", "noun-shield", "noun-bag", "noun-chest", "noun-coin",
		"noun-ring", "adj-shiny", "adj-rusty",
	} {
		id, err := db.CreateLabel(l)
		require.NoError(t, err, "creating label %q", l)
		labels[l] = id
	}

	return svc, db, labels
}

// mustCreate creates and persists an item, returning its ID.
func mustCreate(t *testing.T, svc *item.Service, build func(i *item.Item)) *item.Item {
	t.Helper()

	i := item.New()
	build(i)

	id, err := svc.Create(i)
	require.NoError(t, err)
	require.NotZero(t, id)

	return i
}

func TestNew(t *testing.T) {
	i := item.New()

	assert.NotNil(t, i)
	assert.Equal(t, 100, i.MaxWeight, "default MaxWeight should be 100")
	assert.Zero(t, i.Weight)
	assert.Zero(t, i.ID)
	assert.False(t, i.Container)
	assert.False(t, i.Wearable)
	assert.False(t, i.Worn)
	assert.False(t, i.Created)
}

func TestKindAndAccessors(t *testing.T) {
	i := item.New()

	assert.Equal(t, kind.Item, i.GetKind())

	i.SetID(42)
	assert.Equal(t, uint32(42), i.GetID())

	i.SetLabelID(7)
	assert.Equal(t, uint32(7), i.GetLabelID())
}

func TestServiceCRUD(t *testing.T) {
	svc, _, labels := newTestService(t)

	assert.Zero(t, svc.Count(), "empty service has zero items")

	i := mustCreate(t, svc, func(i *item.Item) {
		i.LabelID = labels["sword"]
		i.NounID = labels["noun-sword"]
		i.Description = "A rusty sword"
	})

	assert.Equal(t, 1, svc.Count())
	assert.NotZero(t, i.GetID())

	i.Description = "A shiny sword"
	require.NoError(t, svc.Update(i))

	got, err := svc.Get().WithID(i.GetID()).First()
	require.NoError(t, err)
	assert.Equal(t, "A shiny sword", got.Description)

	assert.False(t, got.Created)
	require.NoError(t, svc.SetCreated(got, true))

	got, err = svc.Get().WithID(i.GetID()).First()
	require.NoError(t, err)
	assert.True(t, got.Created)
}
