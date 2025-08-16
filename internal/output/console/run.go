package console

import (
	"unicode"

	"github.com/jorgefuertes/thenewquill/internal/log"

	"github.com/gdamore/tcell/v2"
)

// Run starts the console, blocking until the user closes it, you should run it in a goroutine.
func (c *console) Run() {
	for {
		ev := c.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			if c.input.on {
				switch ev.Key() {
				case tcell.KeyEnter:
					c.screen.HideCursor()
					c.input.historyAdd()
					c.input.on = false
					c.Println()
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					if c.input.Len() > 0 && c.input.pos > 0 {
						copy(c.input.current[c.input.pos-1:], c.input.current[c.input.pos:])
						c.input.current = c.input.current[:c.input.Len()-1]
						c.input.pos--
						c.drawInput()

						c.screen.Show()
					}
				case tcell.KeyLeft:
					if c.input.pos > 0 {
						c.input.pos--
						c.moveCursor(-1)
						c.screen.ShowCursor(c.at.Col, c.at.Row)
					}
				case tcell.KeyRight:
					if c.input.pos < c.input.Len() {
						c.input.pos++
						c.moveCursor(1)
						c.screen.ShowCursor(c.at.Col, c.at.Row)
					}
				case tcell.KeyUp, tcell.KeyCtrlR:
					c.input.historyAdd()
					c.clearInputLine()
					c.input.historyPrev()
					c.drawInput()
				case tcell.KeyDown:
					c.input.historyAdd()
					c.clearInputLine()
					c.input.historyNext()
					c.drawInput()
				case tcell.KeyCtrlC, tcell.KeyEsc, tcell.KeyCtrlD:
					c.input.err = ErrCancelledByUser
					c.input.on = false
					c.Println()
					continue
				default:
					if unicode.IsPrint(ev.Rune()) {
						if c.input.Len() >= inputLimit {
							_ = c.screen.Beep()
							continue
						}

						if c.input.pos >= c.input.Len() {
							c.input.current = append(c.input.current, ev.Rune())
							c.input.pos = c.input.Len()
							c.screen.SetContent(c.at.Col, c.at.Row, ev.Rune(), nil, c.style)
						} else {
							c.input.current = append(c.input.current, 0)
							copy(c.input.current[c.input.pos+1:], c.input.current[c.input.pos:])
							c.input.current[c.input.pos] = ev.Rune()
							c.input.pos++
							c.screen.SetContent(c.at.Col, c.at.Row, ev.Rune(), nil, c.style)
						}

						c.moveCursor(1)
					}
				}
			}

			switch ev.Key() {
			case tcell.KeyCtrlC:
				c.Close()
				log.Info("Cancelled by user")

				return
			case tcell.KeyCtrlS:
				c.screen.Sync()
			}
		case *tcell.EventResize:
			c.screen.Sync()
		case *tcell.EventError:
			c.Close()
			log.Error("event error: %s", ev.Error())

			return
		case nil:
			return
		}

		if c.input.on {
			c.screen.ShowCursor(c.at.Col, c.at.Row)
		}
		c.screen.Show()
	}
}

func (c *console) Close() {
	if !c.closed {
		c.screen.SetStyle(tcell.StyleDefault)
		c.screen.Clear()
		c.screen.Clear()
		c.screen.Fini()
		c.closed = true
	}
}
