package console

import "github.com/gdamore/tcell/v2"

func (c *console) scroll() {
	for row := 1; row < c.Rows(); row++ {
		for col := 0; col < c.Cols(); col++ {
			r, _, style, _ := c.screen.GetContent(col, row)
			c.screen.SetContent(col, row-1, r, nil, style)
		}
	}
	for col := 0; col < c.Cols(); col++ {
		c.screen.SetContent(col, c.Rows()-1, ' ', nil, tcell.StyleDefault)
	}

	c.screen.Show()
}
