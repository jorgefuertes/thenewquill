package vars_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/vars"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	st := vars.NewStore()

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
	require.Equal(t, 123, st.GetInt("foo"))
	require.Equal(t, 123.0, st.GetFloat("foo"))
	require.Zero(t, st.GetInt("foo2"))

	// bool
	st.Set("foo", true)
	require.True(t, st.IsSet("foo"))
	require.Equal(t, true, st.GetBool("foo"))
	require.False(t, st.GetBool("foo2"))

	st.Set("foo", 3.14)
	require.True(t, st.IsSet("foo"))
	require.Equal(t, 3.14, st.GetFloat("foo"))
	require.Zero(t, st.GetFloat("foo2"))

	st.Unset("foo")
	require.False(t, st.IsSet("foo"))
}
