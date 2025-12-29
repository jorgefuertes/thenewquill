package variable_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/variable"
	"github.com/jorgefuertes/thenewquill/internal/database"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	db := database.NewDB()
	svc := variable.NewService(db)

	setCases := []struct {
		name           string
		val            any
		wantLabelError bool
		wantSetError   bool
	}{
		{"flashlight.battery", true, false, false},
		{"flashlight.on", true, false, false},
		{"number", 333, false, false},
		{"float", 333.22, false, false},
	}

	getCases := []struct {
		name      string
		val       any
		wantError bool
	}{
		{"flashlight.battery", true, false},
		{"number", 333, false},
		{"float", 333.22, false},
		{"not-found", "", true},
	}

	updateCases := []struct {
		name      string
		val       any
		wantError bool
	}{
		{"number", 666, false},
		{"flashligt.on", true, false},
	}

	for _, tc := range setCases {
		t.Run(tc.name, func(t *testing.T) {
			var id uint32
			var err error

			t.Run("set", func(t *testing.T) {
				id, err = svc.SetByLabel(tc.name, tc.val)
				if tc.wantSetError {
					require.Error(t, err)

					return
				}

				require.NoError(t, err)
				require.NotZero(t, id)
			})

			t.Run("get by id", func(t *testing.T) {
				v, err := svc.Get().WithID(id).First()
				require.NoError(t, err)
				require.NotZero(t, v.LabelID)
				checkValue(t, v, tc.val)
			})

			t.Run("by label", func(t *testing.T) {
				v, err := svc.Get().WithLabel(tc.name).First()
				require.NoError(t, err)
				require.Equal(t, id, v.ID)
				checkValue(t, v, tc.val)
			})
		})
	}

	for _, tc := range getCases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := svc.Get().WithLabel(tc.name).First()
			if tc.wantError {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.NotZero(t, v.ID)
			require.NotZero(t, v.LabelID)

			label, err := db.GetLabel(v.LabelID)
			require.NoError(t, err)
			require.Equal(t, tc.name, label)

			checkValue(t, v, tc.val)
		})
	}

	// UPDATE
	for _, tc := range updateCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("update", func(t *testing.T) {
				id, err := svc.SetByLabel(tc.name, tc.val)
				require.NoError(t, err)
				require.NotZero(t, id)
			})

			t.Run("check", func(t *testing.T) {
				v, err := svc.Get().WithLabel(tc.name).First()
				require.NoError(t, err)

				checkValue(t, v, tc.val)
			})
		})
	}

	// VALIDATE
	t.Run("Validate", func(t *testing.T) {
		err := svc.ValidateAll()
		require.NoError(t, err)
	})
}

func checkValue(t *testing.T, v *variable.Variable, expected any) {
	t.Helper()

	switch val := expected.(type) {
	case int:
		require.Equal(t, val, v.Int())
	case float32, float64:
		require.Equal(t, val, v.Float())
	case bool:
		require.Equal(t, val, v.Bool())
	case string:
		require.Equal(t, expected, v.String())
	default:
		t.Errorf("unexpected type: %T", expected)
	}
}
