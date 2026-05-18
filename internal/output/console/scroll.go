package console

import "github.com/gdamore/tcell/v3"

func (c *console) scroll() {
	for row := 1; row < c.Rows(); row++ {
		for col := 0; col < c.Cols(); col++ {
			str, style, _ := c.screen.Get(col, row)
			c.screen.Put(col, row-1, str, style)
		}
	}
	for col := 0; col < c.Cols(); col++ {
		c.screen.SetContent(col, c.Rows()-1, ' ', nil, tcell.StyleDefault)
	}

	c.screen.Show()
}
