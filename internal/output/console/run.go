package console

import "github.com/gdamore/tcell/v2"

func (c *console) Run() {
	for {
		ev := c.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC:
				c.Close()

				return
			}
		case *tcell.EventResize:
			c.screen.Sync()
		case nil:
			return
		}
	}
}

func (c *console) Close() {
	c.mut.Lock()
	defer c.mut.Unlock()

	c.screen.SetStyle(tcell.StyleDefault)
	c.screen.Clear()
	c.screen.Clear()
	c.screen.Fini()
	c.closed = true
}
