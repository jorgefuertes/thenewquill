package console

import (
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

const runeDelay = time.Millisecond * 5

type console struct {
	mut    sync.Mutex
	screen tcell.Screen
	col    int
	row    int
	style  tcell.Style
	Delay  time.Duration
	closed bool
}

func New() (*console, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	screen.SetCursorStyle(tcell.CursorStyleBlinkingBlock, tcell.ColorWhite)
	_ = screen.Beep()

	return &console{
		mut:    sync.Mutex{},
		screen: screen,
		col:    0,
		row:    0,
		style:  tcell.StyleDefault,
		Delay:  runeDelay,
	}, nil
}
