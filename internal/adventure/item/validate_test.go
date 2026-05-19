package item_test

import (
	"errors"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/item"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemValidate(t *testing.T) {
	testCases := []struct {
		name       string
		build      func(i *item.Item)
		wantOK     bool
		wantErr    error  // sentinel returned by item.Validate (if any)
		wantErrSub string // fallback substring match (for validator messages)
	}{
		{
			name: "valid",
			build: func(i *item.Item) {
				i.LabelID = 1
				i.NounID = 2
				i.Description = "a sword"
			},
			wantOK: true,
		},
		{
			name: "missing LabelID",
			build: func(i *item.Item) {
				i.NounID = 2
				i.Description = "x"
			},
			wantErrSub: "LabelID is required",
		},
		{
			name: "missing NounID",
			build: func(i *item.Item) {
				i.LabelID = 1
				i.Description = "x"
			},
			wantErrSub: "NounID is required",
		},
		{
			name: "missing Description",
			build: func(i *item.Item) {
				i.LabelID = 1
				i.NounID = 2
			},
			wantErrSub: "Description is required",
		},
		{
			name: "weight exceeds max",
			build: func(i *item.Item) {
				i.LabelID = 1
				i.NounID = 2
				i.Description = "heavy"
				i.Weight = 200
				i.MaxWeight = 100
			},
			wantErr: item.ErrWeightShouldBeLessOrEqualThanMaxWeight,
		},
		{
			name: "negative weight",
			build: func(i *item.Item) {
				i.LabelID = 1
				i.NounID = 2
				i.Description = "weird"
				i.Weight = -1
				i.MaxWeight = 100
			},
			wantErr: item.ErrWeightCannotBeNegative,
		},
		{
			name: "negative max weight",
			build: func(i *item.Item) {
				i.LabelID = 1
				i.NounID = 2
				i.Description = "weird"
				i.MaxWeight = -1
			},
			wantErr: item.ErrWeightCannotBeNegative,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			i := item.New()
			tc.build(i)

			err := i.Validate()

			if tc.wantOK {
				require.NoError(t, err)

				return
			}

			require.Error(t, err)

			if tc.wantErr != nil {
				assert.True(t, errors.Is(err, tc.wantErr),
					"got %v, want sentinel %v", err, tc.wantErr)
			}

			if tc.wantErrSub != "" {
				assert.Contains(t, err.Error(), tc.wantErrSub)
			}
		})
	}
}

func TestValidateAll(t *testing.T) {
	t.Run("empty database has no errors", func(t *testing.T) {
		svc, _, _ := newTestService(t)
		assert.Empty(t, svc.ValidateAll())
	})

	t.Run("duplicated noun+adjective is reported", func(t *testing.T) {
		svc, _, labels := newTestService(t)

		mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["sword"]
			i.NounID = labels["noun-sword"]
			i.AdjectiveID = labels["adj-shiny"]
			i.Description = "first"
		})
		mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["shield"]
			i.NounID = labels["noun-sword"]
			i.AdjectiveID = labels["adj-shiny"]
			i.Description = "second"
		})

		errs := svc.ValidateAll()
		require.NotEmpty(t, errs)

		var hit bool
		for _, e := range errs {
			if errors.Is(e, item.ErrDuplicatedNounAdj) {
				hit = true

				break
			}
		}
		assert.True(t, hit, "expected a duplicated noun+adj error")
	})

	t.Run("container overweight is reported", func(t *testing.T) {
		svc, _, labels := newTestService(t)

		bag := mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["bag"]
			i.NounID = labels["noun-bag"]
			i.Description = "small bag"
			i.Container = true
			i.MaxWeight = 5
		})

		mustCreate(t, svc, func(i *item.Item) {
			i.LabelID = labels["sword"]
			i.NounID = labels["noun-sword"]
			i.Description = "heavy stone"
			i.Weight = 10
			i.At = bag.GetID()
		})

		errs := svc.ValidateAll()

		var hit bool
		for _, e := range errs {
			if errors.Is(e, item.ErrContainerCantCarrySoMuch) {
				hit = true

				break
			}
		}
		assert.True(t, hit, "expected a container-overweight error")
	})
}
