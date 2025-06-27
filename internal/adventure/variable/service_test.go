package variable_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/variable"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	d := db.New()
	svc := variable.NewService(d)

	setCases := []struct {
		name           string
		allowDot       db.Allow
		val            any
		wantLabelError bool
		wantSetError   bool
	}{
		{"flashlight.battery", db.AllowDot, true, false, false},
		{"flashlight.on", db.DontAllowDot, true, true, false},
		{"number", db.DontAllowDot, 333, false, false},
		{"array", db.DontAllowDot, []int{1, 2, 3}, false, false},
	}

	updateCases := []struct {
		name      string
		val       any
		wantError bool
	}{
		{"number", 666, false},
	}

	getCases := []struct {
		name      string
		val       any
		wantError bool
	}{
		{"flashlight.battery", true, false},
		{"number", 333, false},
		{"array", []int{1, 2, 3}, false},
		{"not-found", "", true},
	}

	// SET
	for _, tc := range setCases {
		t.Run(tc.name, func(t *testing.T) {
			label, err := d.AddLabel(tc.name, tc.allowDot)
			if tc.wantLabelError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.True(t, label.ID.IsDefined())
				require.NotEmpty(t, label.Name)

				err := svc.Set(label.ID, tc.val)
				if tc.wantSetError {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}

	// GET
	for _, tc := range getCases {
		t.Run(tc.name, func(t *testing.T) {
			label, err := d.GetLabelByName(tc.name)
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.True(t, label.ID.IsDefined())
				require.NotEmpty(t, label.Name)

				v, err := svc.Get(label.ID)
				require.NoError(t, err)
				require.Equal(t, tc.val, v.Value)
			}
		})
	}

	// UPDATE
	for _, tc := range updateCases {
		t.Run(tc.name, func(t *testing.T) {
			label, err := d.GetLabelByName(tc.name)
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.True(t, label.ID.IsDefined())
				require.NotEmpty(t, label.Name)

				err := svc.Set(label.ID, tc.val)
				require.NoError(t, err)

				v, err := svc.Get(label.ID)
				require.NoError(t, err)
				require.Equal(t, tc.val, v.Value)
			}
		})
	}

	// ALL
	t.Run("All", func(t *testing.T) {
		all := svc.All()
		require.NotEmpty(t, all)
		require.NotZero(t, len(all))
	})

	// VALIDATE
	t.Run("Validate", func(t *testing.T) {
		err := svc.ValidateAll()
		require.NoError(t, err)
	})
}
