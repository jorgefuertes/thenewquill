package blob_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/blob"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestService spins up an in-memory database and pre-registers labels.
func newTestService(t *testing.T) (*blob.Service, *database.DB, map[string]uint32) {
	t.Helper()

	db := database.NewDB()
	svc := blob.NewService(db)

	labels := map[string]uint32{}
	for _, l := range []string{"logo", "portrait", "background", "tune"} {
		id, err := db.CreateLabel(l)
		require.NoError(t, err, "creating label %q", l)
		labels[l] = id
	}

	return svc, db, labels
}

// mustCreate persists a blob built by the provided func.
func mustCreate(t *testing.T, svc *blob.Service, build func(b *blob.Blob)) *blob.Blob {
	t.Helper()

	b := blob.New()
	build(b)

	id, err := svc.Create(b)
	require.NoError(t, err)
	require.NotZero(t, id)

	return b
}

func TestNewAndAccessors(t *testing.T) {
	b := blob.New()

	assert.NotNil(t, b)
	assert.Equal(t, kind.Blob, b.GetKind())
	assert.Zero(t, b.GetID())
	assert.Zero(t, b.GetLabelID())
	assert.Nil(t, b.Data, "Data should remain nil until Load is called")

	b.SetID(11)
	assert.Equal(t, uint32(11), b.GetID())

	b.SetLabelID(22)
	assert.Equal(t, uint32(22), b.GetLabelID())
}

func TestLoad(t *testing.T) {
	t.Run("reads file content and infers MIME from extension", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "logo.png")

		// The bytes themselves are arbitrary; the MIME is derived from the .png
		// extension via the stdlib mime package.
		payload := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 'x', 'y'}
		require.NoError(t, os.WriteFile(path, payload, 0o644))

		b := blob.New()
		require.NoError(t, b.Load(path))

		assert.Equal(t, payload, b.Data)
		assert.True(t, strings.HasPrefix(b.Mime, "image/png"),
			"expected image/png MIME, got %q", b.Mime)
	})

	t.Run("reads an empty file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "empty.txt")
		require.NoError(t, os.WriteFile(path, nil, 0o644))

		b := blob.New()
		require.NoError(t, b.Load(path))
		assert.Empty(t, b.Data)
	})

	t.Run("returns error when file does not exist", func(t *testing.T) {
		b := blob.New()
		err := b.Load(filepath.Join(t.TempDir(), "missing.png"))
		require.Error(t, err)
	})

	t.Run("unknown extension yields empty MIME", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "raw.xyz123")
		require.NoError(t, os.WriteFile(path, []byte("data"), 0o644))

		b := blob.New()
		require.NoError(t, b.Load(path))
		assert.Empty(t, b.Mime)
	})
}

func TestServiceCRUD(t *testing.T) {
	svc, _, labels := newTestService(t)

	assert.Zero(t, svc.Count())

	b := mustCreate(t, svc, func(b *blob.Blob) {
		b.LabelID = labels["logo"]
		b.Mime = "image/png"
		b.Data = []byte{1, 2, 3}
	})

	assert.Equal(t, 1, svc.Count())
	assert.NotZero(t, b.GetID())

	b.Data = []byte{9, 9, 9}
	require.NoError(t, svc.Update(b))

	got, err := svc.Get().WithID(b.GetID()).First()
	require.NoError(t, err)
	assert.Equal(t, []byte{9, 9, 9}, got.Data)
}

func TestQuery(t *testing.T) {
	svc, _, labels := newTestService(t)

	logo := mustCreate(t, svc, func(b *blob.Blob) {
		b.LabelID = labels["logo"]
		b.Mime = "image/png"
		b.Data = []byte{1}
	})
	portrait := mustCreate(t, svc, func(b *blob.Blob) {
		b.LabelID = labels["portrait"]
		b.Mime = "image/jpeg"
		b.Data = []byte{2}
	})
	tune := mustCreate(t, svc, func(b *blob.Blob) {
		b.LabelID = labels["tune"]
		b.Mime = "audio/mpeg"
		b.Data = []byte{3}
	})

	t.Run("Count", func(t *testing.T) {
		assert.Equal(t, 3, svc.Get().Count())
	})

	t.Run("WithID", func(t *testing.T) {
		got, err := svc.Get().WithID(portrait.GetID()).First()
		require.NoError(t, err)
		assert.Equal(t, "image/jpeg", got.Mime)
	})

	t.Run("WithLabel", func(t *testing.T) {
		got, err := svc.Get().WithLabel("logo").First()
		require.NoError(t, err)
		assert.Equal(t, logo.GetID(), got.GetID())
	})

	t.Run("WithLabelID", func(t *testing.T) {
		got, err := svc.Get().WithLabelID(labels["tune"]).First()
		require.NoError(t, err)
		assert.Equal(t, tune.GetID(), got.GetID())
	})

	t.Run("WithNoID excludes", func(t *testing.T) {
		assert.Equal(t, 2, svc.Get().WithNoID(logo.GetID()).Count())
	})

	t.Run("Exists", func(t *testing.T) {
		assert.True(t, svc.Get().WithLabel("logo").Exists())
		assert.False(t, svc.Get().WithLabel("nonexistent").Exists())
	})

	t.Run("First missing returns error", func(t *testing.T) {
		_, err := svc.Get().WithID(99999).First()
		require.Error(t, err)
	})
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		name       string
		build      func(b *blob.Blob)
		wantOK     bool
		wantErrSub string
	}{
		{
			name: "valid",
			build: func(b *blob.Blob) {
				b.LabelID = 1
				b.Mime = "image/png"
				b.Data = []byte{1, 2, 3}
			},
			wantOK: true,
		},
		{
			name: "missing LabelID",
			build: func(b *blob.Blob) {
				b.Mime = "image/png"
				b.Data = []byte{1}
			},
			wantErrSub: "LabelID is required",
		},
		{
			name: "missing Mime",
			build: func(b *blob.Blob) {
				b.LabelID = 1
				b.Data = []byte{1}
			},
			wantErrSub: "Mime is required",
		},
		{
			name: "missing Data",
			build: func(b *blob.Blob) {
				b.LabelID = 1
				b.Mime = "image/png"
			},
			wantErrSub: "Data is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := blob.New()
			tc.build(b)

			err := b.Validate()

			if tc.wantOK {
				require.NoError(t, err)

				return
			}

			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.wantErrSub)
		})
	}
}

func TestValidateAll(t *testing.T) {
	t.Run("empty database has no errors", func(t *testing.T) {
		svc, _, _ := newTestService(t)
		assert.Empty(t, svc.ValidateAll())
	})

	t.Run("invalid blob is reported with context", func(t *testing.T) {
		svc, _, labels := newTestService(t)

		// Valid blob.
		mustCreate(t, svc, func(b *blob.Blob) {
			b.LabelID = labels["logo"]
			b.Mime = "image/png"
			b.Data = []byte{1}
		})

		// Invalid: missing Mime.
		mustCreate(t, svc, func(b *blob.Blob) {
			b.LabelID = labels["portrait"]
			b.Data = []byte{2}
		})

		errs := svc.ValidateAll()
		require.Len(t, errs, 1)
		assert.Contains(t, errs[0].Error(), "Mime is required")
		assert.Contains(t, errs[0].Error(), "portrait",
			"error should mention the offending blob's label")
	})
}
