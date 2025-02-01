package store_test

import (
	"testing"

	"thenewquill/internal/adventure/store"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	st := store.New()

	// string
	st.Set("foo", "bar")
	require.True(t, st.IsSet("foo"))
	require.Equal(t, "bar", st.Get("foo"))
	require.Empty(t, st.Get("foo2"))
	require.Equal(t, "", st.Get("foo2"))

	// unset
	st.Unset("foo")
	require.False(t, st.IsSet("foo"))
	require.Empty(t, st.Get("foo"))

	// number
	st.Set("foo", 123)
	require.True(t, st.IsSet("foo"))
	require.Equal(t, 123.0, st.GetNumber("foo"))
	require.Greater(t, st.GetNumber("foo"), float64(20))
	require.Zero(t, st.GetNumber("foo2"))

	// bool
	st.Set("foo", true)
	require.True(t, st.IsSet("foo"))
	require.Equal(t, true, st.GetBool("foo"))
	require.False(t, st.GetBool("foo2"))

	st.Set("foo", 3.14)
	require.True(t, st.IsSet("foo"))
	require.Equal(t, 3.14, st.GetNumber("foo"))
	require.Zero(t, st.GetNumber("foo2"))

	st.Unset("foo")
	require.False(t, st.IsSet("foo"))
}
