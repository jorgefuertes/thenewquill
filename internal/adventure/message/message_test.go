package message_test

import (
	"testing"

	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/adventure/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	noCoinsInTheBox = "There's no coins in the box."
	oneCoinInTheBox = "There's one coin in the box."
	helloWorld      = "Hello, World!"
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
				Text: helloWorld,
			},
			calls: []msgCall{
				{args: []any{}, want: helloWorld},
				{args: []any{0, 1}, want: helloWorld},
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
				Text: noCoinsInTheBox,
				Plurals: [2]string{
					oneCoinInTheBox,
					"There's _ coins in the box.",
				},
			},
			calls: []msgCall{
				{args: []any{0}, want: noCoinsInTheBox},
				{args: []any{1}, want: oneCoinInTheBox},
				{args: []any{2}, want: "There's 2 coins in the box."},
				{args: []any{18}, want: "There's 18 coins in the box."},
				{args: []any{1.5}, want: "There's 1.50 coins in the box."},
				{args: []any{"zero"}, want: noCoinsInTheBox},
				{args: []any{"one"}, want: oneCoinInTheBox},
				{args: []any{"una"}, want: oneCoinInTheBox},
				{args: []any{"1"}, want: oneCoinInTheBox},
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
	m := message.New()
	m.Text = helloWorld
	require.NotEmpty(t, m)
	assert.Equal(t, helloWorld, m.String())

	m.Text = "No coins."
	m.SetPlurals([2]string{"One coin.", "Many coins."})
	assert.Equal(t, "No coins.", m.String())
	assert.Equal(t, "One coin.", m.Stringf(1))
	assert.Equal(t, "Many coins.", m.Stringf(2))

	m.ID = 5
	assert.Equal(t, uint32(5), m.ID)

	m.ID = 10
	assert.Equal(t, uint32(10), m.ID)

	assert.Equal(t, kind.Message, kind.KindOf(m))

	m = message.New()
	m.Text = "No one"
	m.SetPlural(message.One, "Only one")
	m.SetPlural(message.Many, "Many")
	assert.Equal(t, "No one", m.String())
	assert.Equal(t, "Only one", m.Stringf(1))
	assert.Equal(t, "Many", m.Stringf(2))
}

func TestValidate(t *testing.T) {
	m := message.New()
	require.Error(t, m.Validate())
	assert.ErrorContains(t, m.Validate(), "LabelID is required")
	m.LabelID = 7
	assert.ErrorContains(t, m.Validate(), "Text is required")

	m.Text = helloWorld
	m.Plurals[message.One] = "One"
	require.Error(t, m.Validate())
	assert.ErrorIs(t, message.ErrUndefinedPlural, m.Validate())
	m.Plurals[message.Many] = "Many"
	require.NoError(t, m.Validate())
}
