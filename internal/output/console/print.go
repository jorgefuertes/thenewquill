package console

import (
	"fmt"
	"time"

	"github.com/mitchellh/go-wordwrap"
)

func (c *console) Print(s string) {
	c.mut.Lock()
	defer c.mut.Unlock()

	beforeCol := c.col
	beforeRow := c.row

	for _, r := range s {
		switch r {
		case '\n':
			c.row++
			c.col = 0
		case '\r':
			c.col = 0
		case '\t':
			c.col += 4
			c.col -= c.col % 4
		default:
			c.screen.SetContent(c.col, c.row, r, nil, c.style.Reverse(true))
			if c.Delay > 0 {
				c.screen.Show()
				time.Sleep(c.Delay)
			}

			c.screen.SetContent(c.col, c.row, r, nil, c.style)
			c.screen.Show()
			c.col++
		}

		if c.col >= c.Cols() {
			c.row++
			c.col = 0
		}

		if c.row >= c.Rows() {
			c.row = c.Rows() - 1
			c.scroll()
		}
	}

	if c.Delay > 0 && beforeCol != c.col || beforeRow != c.row {
		c.screen.Show()
	}
}

func (c *console) Println(a ...any) {
	c.Print(fmt.Sprintln(a...))
}

func (c *console) Printf(format string, a ...any) {
	c.Print(fmt.Sprintf(format, a...))
}

func (c *console) Cls() {
	c.mut.Lock()
	defer c.mut.Unlock()

	c.screen.Clear()
	c.row = 0
	c.col = 0
	c.screen.Show()
}

func (c *console) WrapPrint(text string) {
	if c.Cols() <= 0 {
		return
	}

	if len(text) <= c.Cols() {
		c.Print(text)

		return
	}

	c.Print(wordwrap.WrapString(text, uint(c.Cols())))
}

func (c *console) WrapPrintf(format string, a ...any) {
	c.WrapPrint(fmt.Sprintf(format, a...))
}
