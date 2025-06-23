package message_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessages(t *testing.T) {
	type msgCall struct {
		args []any
		want string
	}

	type testCase struct {
		name  string
		msg   message.Message
		calls []msgCall
	}

	testCases := []testCase{
		{
			name: "plain",
			msg: message.Message{
				Text: "Hello, World!",
			},
			calls: []msgCall{
				{args: []any{}, want: "Hello, World!"},
				{args: []any{0, 1}, want: "Hello, World!"},
			},
		},
		{
			name: "substitution",
			msg: message.Message{
				Text: "Hello _. How are you?",
			},
			calls: []msgCall{
				{args: []any{}, want: "Hello ?. How are you?"},
				{args: []any{"Test"}, want: "Hello Test. How are you?"},
				{args: []any{"Test", "Another Arg", 2}, want: "Hello Test. How are you?"},
				{args: []any{"You", "Another Arg", 2, true}, want: "Hello You. How are you?"},
			},
		},
		{
			name: "substitution2",
			msg: message.Message{
				Text: "The arg one is _ and the arg two is _.",
			},
			calls: []msgCall{
				{args: []any{1, 2}, want: "The arg one is 1 and the arg two is 2."},
				{args: []any{}, want: "The arg one is ? and the arg two is ?."},
				{args: []any{1}, want: "The arg one is 1 and the arg two is ?."},
			},
		},
		{
			name: "pluralization",
			msg: message.Message{
				Text: "There's no coins in the box.",
				Plurals: [2]string{
					"There's one coin in the box.",
					"There's _ coins in the box.",
				},
			},
			calls: []msgCall{
				{args: []any{0}, want: "There's no coins in the box."},
				{args: []any{1}, want: "There's one coin in the box."},
				{args: []any{2}, want: "There's 2 coins in the box."},
				{args: []any{18}, want: "There's 18 coins in the box."},
				{args: []any{1.5}, want: "There's 1.50 coins in the box."},
				{args: []any{"zero"}, want: "There's no coins in the box."},
				{args: []any{"one"}, want: "There's one coin in the box."},
				{args: []any{"una"}, want: "There's one coin in the box."},
				{args: []any{"1"}, want: "There's one coin in the box."},
				{args: []any{"a lot of"}, want: "There's a lot of coins in the box."},
				{args: []any{true}, want: "There's true coins in the box."},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, c := range tc.calls {
				got := tc.msg.Stringf(c.args...)
				if got != c.want {
					t.Errorf("got %q, want %q", got, c.want)
				}
			}
		})
	}
}

func TestStoreable(t *testing.T) {
	m := message.New("Hello, World!")
	require.NotEmpty(t, m)
	assert.Equal(t, "Hello, World!", m.String())

	m = message.New("No coins.")
	m.SetPlurals([2]string{"One coin.", "Many coins."})
	assert.Equal(t, "No coins.", m.String())
	assert.Equal(t, "One coin.", m.Stringf(1))
	assert.Equal(t, "Many coins.", m.Stringf(2))

	m.ID = db.ID(5)
	assert.Equal(t, db.ID(5), m.GetID())

	s := m.SetID(db.ID(10))
	assert.Equal(t, db.ID(10), s.GetID())

	assert.Equal(t, db.Messages, m.GetKind())

	m = message.New("No one")
	m.SetPlural(message.One, "Only one")
	m.SetPlural(message.Many, "Many")
	assert.Equal(t, "No one", m.String())
	assert.Equal(t, "Only one", m.Stringf(1))
	assert.Equal(t, "Many", m.Stringf(2))
}

func TestValidate(t *testing.T) {
	m := message.New("Hello, World!")
	require.Error(t, m.Validate(db.DontAllowNoID))
	assert.ErrorIs(t, db.ErrUndefinedLabel, m.Validate(db.DontAllowNoID))

	m.ID = db.ID(3)
	require.Error(t, m.Validate(db.DontAllowNoID))
	assert.ErrorIs(t, db.ErrInvalidLabelID, m.Validate(db.DontAllowNoID))

	m.ID = db.ID(4)
	require.NoError(t, m.Validate(db.DontAllowNoID))

	m.Text = ""
	require.Error(t, m.Validate(db.DontAllowNoID))
	assert.ErrorIs(t, message.ErrUndefinedText, m.Validate(db.DontAllowNoID))

	m.Text = "Hello, World!"
	m.Plurals[message.One] = "One"
	require.Error(t, m.Validate(db.DontAllowNoID))
	assert.ErrorIs(t, message.ErrUndefinedPlural, m.Validate(db.DontAllowNoID))
	m.Plurals[message.Many] = "Many"
	require.NoError(t, m.Validate(db.DontAllowNoID))
}
