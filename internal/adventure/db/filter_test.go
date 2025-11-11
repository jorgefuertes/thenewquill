package db

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adapter"
	"github.com/jorgefuertes/thenewquill/internal/adventure/id"
	"github.com/stretchr/testify/require"
)

type testStoreable struct {
	ID      id.ID
	Title   string
	At      id.ID
	OK      bool
	NOOK    bool
	Weight  int16
	Names   []string
	Numbers []int
}

var _ adapter.Storeable = testStoreable{}

func (s testStoreable) GetID() id.ID {
	return s.ID
}

func (s testStoreable) SetID(id id.ID) adapter.Storeable {
	s.ID = id

	return s
}

func (s testStoreable) Validate(allowNoID bool) error {
	return nil
}

func TestMatches(t *testing.T) {
	type testCase struct {
		name     string
		sut      *testStoreable
		filters  []filter
		expected bool
	}

	testItem := &testStoreable{
		ID:      id.ID(7),
		Title:   "This is just a Test Title",
		At:      id.ID(10),
		OK:      true,
		NOOK:    false,
		Weight:  10,
		Names:   []string{"test1", "test2", "test3"},
		Numbers: []int{1, 2, 45, 67},
	}

	testCases := []testCase{
		{
			name: "Equal id.ID",
			sut:  testItem,
			filters: []filter{
				{Equal, "ID", id.ID(7)},
			},
			expected: true,
		},
		{
			name: "Equal id.ID by number",
			sut:  testItem,
			filters: []filter{
				{Equal, "ID", 7},
			},
			expected: true,
		},
		{
			name: "Equal id.ID false",
			sut:  testItem,
			filters: []filter{
				{Equal, "ID", id.ID(23)},
			},
			expected: false,
		},
		{
			name: "NotEqual id.ID",
			sut:  testItem,
			filters: []filter{
				{NotEqual, "ID", id.ID(23)},
			},
			expected: true,
		},
		{
			name: "Equal weight",
			sut:  testItem,
			filters: []filter{
				{Equal, "weight", 10},
			},
			expected: true,
		},
		{
			name: "Equal weight false",
			sut:  testItem,
			filters: []filter{
				{Equal, "Weight", 15},
			},
			expected: false,
		},
		{
			name: "NotEqual weight",
			sut:  testItem,
			filters: []filter{
				{NotEqual, "Weight", 15},
			},
			expected: true,
		},
		{
			name: "Equal OK",
			sut:  testItem,
			filters: []filter{
				{Equal, "ok", true},
			},
			expected: true,
		},
		{
			name: "Equal NOOK",
			sut:  testItem,
			filters: []filter{
				{Equal, "NOok", false},
			},
			expected: true,
		},
		{
			name: "Equal Title",
			sut:  testItem,
			filters: []filter{
				{Equal, "Title", "This is just a Test Title"},
			},
			expected: true,
		},
		{
			name: "Contains Title",
			sut:  testItem,
			filters: []filter{
				{Contains, "Title", "Test Title"},
			},
			expected: true,
		},
		{
			name: "Contains Title false",
			sut:  testItem,
			filters: []filter{
				{Contains, "Title", "Nothing"},
			},
			expected: false,
		},
		{
			name: "NotContains Title",
			sut:  testItem,
			filters: []filter{
				{NotContains, "Title", "Nothing"},
			},
			expected: true,
		},
		{
			name: "Contains Names",
			sut:  testItem,
			filters: []filter{
				{Contains, "Names", "test2"},
			},
			expected: true,
		},
		{
			name: "Contains Names false",
			sut:  testItem,
			filters: []filter{
				{Contains, "Names", "not-found"},
			},
			expected: false,
		},
		{
			name: "NotContains Names",
			sut:  testItem,
			filters: []filter{
				{NotContains, "Names", "not-found"},
			},
			expected: true,
		},
		{
			name: "NotContains Names false",
			sut:  testItem,
			filters: []filter{
				{NotContains, "Names", "test1"},
			},
			expected: false,
		},
		{
			name: "Contains Numbers",
			sut:  testItem,
			filters: []filter{
				{Contains, "nUmbers", 45},
			},
			expected: true,
		},
		{
			name: "Contains Numbers false",
			sut:  testItem,
			filters: []filter{
				{Contains, "numbers", 46},
			},
			expected: false,
		},
		{
			name: "NotContains Numbers",
			sut:  testItem,
			filters: []filter{
				{NotContains, "Numbers", 46},
			},
			expected: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := matches(tt.sut, tt.filters...)
			require.Equal(t, tt.expected, result, "Result is %t, expecting %t", result, tt.expected)
		})
	}
}
