package console

import (
	"fmt"
	"time"

	"github.com/mitchellh/go-wordwrap"
)

func (c *console) print(r rune) {
	switch r {
	case '\n':
		c.incRow()
	case '\r':
		c.at.Col = 0
	case '\t':
		c.at.Col += 4
		c.at.Col -= c.at.Col % 4
	default:
		c.screen.SetContent(c.at.Col, c.at.Row, r, nil, c.style)
		c.moveCursor(1)
	}
}

func (c *console) Print(s string) {
	c.mut.Lock()
	defer c.mut.Unlock()

	for _, r := range s {
		c.print(r)
		if c.Delay > 0 {
			c.screen.SetContent(c.at.Col, c.at.Row, cursorShape, nil, c.style)
			c.screen.Show()
			time.Sleep(c.Delay)
			c.screen.SetContent(c.at.Col, c.at.Row, ' ', nil, c.style)
		}

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
	c.at.Row = 0
	c.at.Col = 0
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
