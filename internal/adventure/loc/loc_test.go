package loc_test

import (
	"testing"

	"thenewquill/internal/adventure/loc"
	"thenewquill/internal/adventure/voc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocations(t *testing.T) {
	locs := loc.New()

	north := &voc.Word{Label: "north", Type: voc.Verb}
	east := &voc.Word{Label: "east", Type: voc.Verb}
	west := &voc.Word{Label: "west", Type: voc.Verb}
	south := &voc.Word{Label: "shouth", Type: voc.Verb}

	definitions := []struct {
		label string
		title string
		desc  string
		exits map[*voc.Word]string
	}{
		{"loc-001", "loc 001 title", "loc 001 desc", map[*voc.Word]string{east: "loc-002"}},
		{"loc-002", "loc 002 title", "loc 002 desc", map[*voc.Word]string{west: "loc-001"}},
		{"loc-003", "loc 003 title", "loc 003 desc", map[*voc.Word]string{north: "loc-004", south: "loc-005"}},
		{"loc-004", "loc 004 title", "loc 004 desc", map[*voc.Word]string{east: "loc-003"}},
		{"loc-005", "loc 005 title", "loc 005 desc", map[*voc.Word]string{west: "loc-003"}},
		{"loc-006", "loc 006 title", "loc 006 desc", map[*voc.Word]string{north: "loc-007", south: "loc-008"}},
		{"loc-007", "loc 007 title", "loc 007 desc", map[*voc.Word]string{east: "loc-006"}},
		{"loc-008", "loc 008 title", "loc 008 desc", map[*voc.Word]string{west: "loc-006"}},
		{"loc-009", "loc 009 title", "loc 009 desc", map[*voc.Word]string{north: "loc-001", south: "loc-002"}},
	}

	t.Run("create locations", func(t *testing.T) {
		for _, d := range definitions {
			l := locs.Set(d.label, d.title, d.desc)
			require.NotNil(t, l, "set should return a location for %s", d.label)
			assert.Equal(t, d.title, l.Title, "location %s title should match", d.label)
			assert.Equal(t, d.desc, l.Description, "location %s description should match", d.label)
			assert.Len(t, l.Conns, 0, "conns sould be empty for %s", d.label)
			for w, label := range d.exits {
				to := locs.Get(label)
				if to == nil {
					// create an empty location
					to = locs.Set(label, "", "")
				}

				l.SetConn(w, to)
			}
		}
	})

	t.Run("exists", func(t *testing.T) {
		for _, d := range definitions {
			assert.True(t, locs.Exists(d.label), "location %s should exist", d.label)
		}
	})

	t.Run("get", func(t *testing.T) {
		for _, d := range definitions {
			l := locs.Get(d.label)
			assert.NotNil(t, l, "get should return a location for %s", d.label)
			assert.Equal(t, d.label, l.Label, "location %s label should match", d.label)
			assert.Equal(t, d.title, l.Title, "location %s title should match", d.label)
			assert.Equal(t, d.desc, l.Description, "location %s description should match", d.label)
			for w, label := range d.exits {
				dest := l.GetConn(w)
				require.NotNil(t, dest, "word %s should have a destination", w)
				assert.Equal(t, label, dest.Label, "word %s destination should match", w)
			}
		}
	})
}
