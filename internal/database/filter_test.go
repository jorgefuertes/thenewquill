package database

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/stretchr/testify/require"
)

func TestMatches(t *testing.T) {
	type TestItem struct {
		ID      uint32
		LabelID uint32
		Title   string
		At      uint32
		OK      bool
		NOOK    bool
		Weight  int
		Names   []string
		Numbers []int
	}

	type testCase struct {
		name     string
		sut      *TestItem
		filters  []filter
		expected bool
	}

	const testItemTile = "This is just a Test Title"

	testItem1 := &TestItem{
		ID:      7,
		LabelID: 3,
		Title:   testItemTile,
		At:      2,
		OK:      true,
		NOOK:    false,
		Weight:  10,
		Names:   []string{"one", "two", "three"},
		Numbers: []int{1, 2, 3},
	}

	db := NewDB()
	id, err := db.CreateLabel("test-label")
	require.NoError(t, err)
	testItem1.LabelID = id

	testCases := []testCase{
		{
			name:     "Equal ID",
			sut:      testItem1,
			filters:  []filter{FilterByID(7)},
			expected: true,
		},
		{
			name:     "Equal LabelID",
			sut:      testItem1,
			filters:  []filter{FilterByLabelID(3)},
			expected: true,
		},
		{
			name:     "Equal ID false",
			sut:      testItem1,
			filters:  []filter{FilterByID(23)},
			expected: false,
		},
		{
			name:     "NotEqual ID",
			sut:      testItem1,
			filters:  []filter{{NotEqual, "ID", 23}},
			expected: true,
		},
		{
			name: "Equal weight",
			sut:  testItem1,
			filters: []filter{
				{Equal, "weight", 10},
			},
			expected: true,
		},
		{
			name: "Equal weight false",
			sut:  testItem1,
			filters: []filter{
				{Equal, "Weight", 15},
			},
			expected: false,
		},
		{
			name: "NotEqual weight",
			sut:  testItem1,
			filters: []filter{
				{NotEqual, "Weight", 15},
			},
			expected: true,
		},
		{
			name: "Equal OK",
			sut:  testItem1,
			filters: []filter{
				{Equal, "ok", true},
			},
			expected: true,
		},
		{
			name: "Equal NOOK",
			sut:  testItem1,
			filters: []filter{
				{Equal, "NOok", false},
			},
			expected: true,
		},
		{
			name: "Equal Title",
			sut:  testItem1,
			filters: []filter{
				{Equal, "Title", "This is just a Test Title"},
			},
			expected: true,
		},
		{
			name: "Contains in Title",
			sut:  testItem1,
			filters: []filter{
				{Contains, "Title", "Test Title"},
			},
			expected: true,
		},
		{
			name: "Contains in Title false",
			sut:  testItem1,
			filters: []filter{
				{Contains, "Title", "Nothing"},
			},
			expected: false,
		},
		{
			name: "NotContains in Title",
			sut:  testItem1,
			filters: []filter{
				{NotContains, "Title", "Nothing"},
			},
			expected: true,
		},
		{
			name: "Contains in Names",
			sut:  testItem1,
			filters: []filter{
				{Contains, "Names", "one"},
			},
			expected: true,
		},
		{
			name: "Contains in Names false",
			sut:  testItem1,
			filters: []filter{
				{Contains, "Names", "four"},
			},
			expected: false,
		},
		{
			name: "NotContains in Names",
			sut:  testItem1,
			filters: []filter{
				{NotContains, "Names", "five"},
			},
			expected: true,
		},
		{
			name: "NotContains in Names false",
			sut:  testItem1,
			filters: []filter{
				{NotContains, "Names", "two"},
			},
			expected: false,
		},
		{
			name: "Contains in Numbers",
			sut:  testItem1,
			filters: []filter{
				{Contains, "nUmbers", 2},
			},
			expected: true,
		},
		{
			name: "Contains in Numbers false",
			sut:  testItem1,
			filters: []filter{
				{Contains, "numbers", 4},
			},
			expected: false,
		},
		{
			name: "NotContains in Numbers",
			sut:  testItem1,
			filters: []filter{
				{NotContains, "Numbers", 5},
			},
			expected: true,
		},
		{
			name:     "By Label",
			sut:      testItem1,
			filters:  []filter{FilterByLabel("test-label")},
			expected: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := Record{
				LabelID: tt.sut.LabelID,
				Kind:    kind.Item,
				Data:    []byte{},
			}

			err := r.Marshal(tt.sut)
			require.NoError(t, err)

			result := db.matchesAllFilters(r, tt.filters...)
			require.Equal(t, tt.expected, result, "Result is %t, expecting %t for: %v", result, tt.expected, tt.filters)
		})
	}
}
