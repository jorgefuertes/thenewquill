package console

import "github.com/gdamore/tcell/v2"

func (c *console) SetStyle(style tcell.Style) {
	c.style = style
}
