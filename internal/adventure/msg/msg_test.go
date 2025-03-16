package msg_test

import (
	"testing"

	"thenewquill/internal/adventure/msg"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	m := msg.New("foo", "bar")
	assert.Equal(t, "foo", m.Label)
	assert.Equal(t, "bar", m.Text)
	assert.Equal(t, "bar", m.String())

	m = msg.New("foo", "This is a _ message.")
	assert.Equal(t, "This is a test message.", m.Stringf("test"))
	assert.Equal(t, "This is a 1 message.", m.Stringf(1))
	assert.Equal(t, "This is a true message.", m.Stringf(true))

	m = msg.New("foo.zero", "There are no foos.")
	assert.Equal(t, "", m.Text)
	m.SetPlurals([3]string{"There are no foos.", "There is one foo.", "There are _ foos."})
	assert.Equal(t, "There are no foos.", m.Stringf(0))
	assert.Equal(t, "There is one foo.", m.Stringf(1))
	assert.Equal(t, "There are 2 foos.", m.Stringf(2))
	assert.Equal(t, "There are 34 foos.", m.Stringf(34))
	assert.Equal(t, "There are true foos.", m.Stringf(true))
	assert.Equal(t, "There are no foos.", m.Stringf("zero"))
	assert.Equal(t, "There are no foos.", m.Stringf("cero"))
	assert.Equal(t, "There are no foos.", m.Stringf("0"))
	assert.Equal(t, "There are no foos.", m.Stringf(0.0))
	assert.Equal(t, "There is one foo.", m.Stringf("one"))
	assert.Equal(t, "There is one foo.", m.Stringf("un"))
	assert.Equal(t, "There is one foo.", m.Stringf("uno"))
	assert.Equal(t, "There is one foo.", m.Stringf("una"))
	assert.Equal(t, "There is one foo.", m.Stringf("1"))
	assert.Equal(t, "There is one foo.", m.Stringf(1.0))
	assert.Equal(t, "There are many foos.", m.Stringf("many"))
	assert.Equal(t, "There are several foos.", m.Stringf("several"))
	assert.Equal(t, "There are 98.05 foos.", m.Stringf(98.05))
}

func TestIsPluralized(t *testing.T) {
	m := msg.New("foo", "bar")
	assert.False(t, m.IsPluralized())

	m = msg.New("foo.zero", "There's no foos")
	assert.True(t, m.IsPluralized())
}
