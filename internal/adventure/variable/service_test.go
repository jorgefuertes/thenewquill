package variable_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/database"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database/primitive"
	"github.com/jorgefuertes/thenewquill/internal/adventure/variable"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	db := database.New()
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
			var id primitive.ID
			var err error

			t.Run("set", func(t *testing.T) {
				id, err = svc.CreateWithLabel(tc.name, tc.val)
				if tc.wantSetError {
					require.Error(t, err)

					return
				}

				require.NoError(t, err)
				require.True(t, id.IsDefinedID())
			})

			t.Run("get by id", func(t *testing.T) {
				v, err := svc.Get(id)
				require.NoError(t, err)
				require.True(t, v.LabelID.IsDefinedID())
				require.Equal(t, tc.val, v.Value)
			})

			t.Run("by label", func(t *testing.T) {
				v, err := svc.GetByLabel(tc.name)
				require.NoError(t, err)
				require.Equal(t, id, v.ID)
				require.Equal(t, tc.val, v.Value)
			})
		})
	}

	for _, tc := range getCases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := svc.GetByLabel(tc.name)
			if tc.wantError {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.True(t, v.ID.IsDefinedID())
			require.True(t, v.LabelID.IsDefinedID())

			label, err := db.GetLabel(v.LabelID)
			require.NoError(t, err)
			require.Equal(t, tc.name, label)

			switch val := tc.val.(type) {
			case int:
				require.Equal(t, val, v.Int())
			case float32, float64:
				require.Equal(t, val, v.Float())
			case bool:
				require.Equal(t, val, v.Bool())
			case string:
				require.Equal(t, tc.val, v.String())
			default:
				t.Errorf("unexpected type: %T", tc.val)
			}
		})
	}

	// UPDATE
	for _, tc := range updateCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("update", func(t *testing.T) {
				v, err := svc.GetByLabel(primitive.Label(tc.name))
				require.NoError(t, err)

				v.Set(tc.val)
			})

			t.Run("check", func(t *testing.T) {
				tcVar := &variable.Variable{}
				tcVar.Set(tc.val)

				v, err := svc.GetByLabel(primitive.Label(tc.name))
				require.NoError(t, err)

				require.Equal(t, tcVar.Value, v.Value)
			})
		})
	}

	// VALIDATE
	t.Run("Validate", func(t *testing.T) {
		err := svc.ValidateAll()
		require.NoError(t, err)
	})
}
