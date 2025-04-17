package console

import (
	"time"
)

type input struct {
	prompt  string
	timeout time.Duration
	on      bool
	current []rune
	history [][]rune
	index   int
	pos     int
	orig    At
	err     error
}

func newInput() input {
	return input{
		prompt:  "",
		timeout: 0,
		on:      false,
		current: []rune{},
		history: [][]rune{},
		index:   0,
		pos:     0,
		orig: At{
			Row: 0,
			Col: 0,
		},
		err: nil,
	}
}

func (i input) Len() int {
	return len(i.current)
}

func (i input) String() string {
	return string(i.current)
}

func (i *input) reset(prompt string) {
	i.on = false
	i.prompt = prompt
	i.current = []rune{}
	i.index = len(i.history) - 1
	i.pos = 0
	i.orig = At{
		Row: 0,
		Col: 0,
	}
	i.err = nil
}

// clearInputLine clears the input buffer and display
func (c *console) clearInputLine() {
	c.at = c.input.orig
	c.moveCursor(len(c.input.prompt) + 1)
	for i := 0; i < c.input.Len(); i++ {
		c.print(' ')
	}

	c.screen.Show()
}

func (c *console) Input(prompt string, timeout time.Duration) (string, error) {
	c.input.reset(prompt)
	c.input.prompt = prompt
	limit := time.Now().Add(timeout)

	if c.at.Col > 0 {
		c.Println()
	}
	c.input.orig = c.at

	c.drawInput()
	c.input.on = true

	for c.input.on {
		time.Sleep(time.Millisecond * 50)
		if timeout > 0 && time.Now().After(limit) {
			return string(c.input.current), ErrTimedOut
		}
	}

	return string(c.input.current), c.input.err
}

func (c *console) drawInput() {
	c.at = c.input.orig

	for _, r := range c.input.prompt {
		c.print(r)
	}
	c.print(' ')
	for _, r := range c.input.current {
		c.print(r)
	}
	c.print(' ')

	c.at = c.input.orig
	c.moveCursor(len(c.input.prompt) + 1 + c.input.pos)

	c.screen.ShowCursor(c.at.Col, c.at.Row)
	c.screen.Show()
}
