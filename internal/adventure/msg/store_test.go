package msg_test

import (
	"testing"

	"thenewquill/internal/adventure/msg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgStore(t *testing.T) {
	// I can create an empty store
	s := msg.NewStore()
	assert.Equal(t, 0, s.Len())

	// I can add a message
	err := s.Set(msg.New(msg.SystemMsg, "foo", "bar"))
	require.NoError(t, err)
	assert.Equal(t, 1, s.Len())

	// I can add a message with plurals
	err = s.Set(msg.New(msg.SystemMsg, "foos.zero", "There are no foos."))
	require.NoError(t, err)
	assert.Equal(t, 2, s.Len())

	// I can add the remaining plurals
	err = s.Set(msg.New(msg.SystemMsg, "foos.one", "There is one foo."))
	require.NoError(t, err)
	assert.Equal(t, 2, s.Len())
	err = s.Set(msg.New(msg.SystemMsg, "foos.many", "There are _ foos."))
	require.NoError(t, err)
	assert.Equal(t, 2, s.Len())

	// I can add another message
	err = s.Set(msg.New(msg.SystemMsg, "bar", "There is a bar."))
	require.NoError(t, err)
	assert.Equal(t, 3, s.Len())

	// I can get a plural string from a message
	m := s.Get(msg.SystemMsg, "foos")
	require.NotNil(t, m)
	require.Equal(t, "There are 7 foos.", m.Stringf(7))
}
