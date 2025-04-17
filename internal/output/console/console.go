package console

import (
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	runeDelay   = time.Millisecond * 5
	historySize = 100
	cursorShape = '█'
	inputLimit  = 256
)

type At struct {
	Row int
	Col int
}

type console struct {
	mut     *sync.Mutex
	screen  tcell.Screen
	at      At
	atStack []At
	style   tcell.Style
	Delay   time.Duration
	closed  bool
	input   input
}

func New() (*console, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	screen.HideCursor()
	_ = screen.Beep()

	return &console{
		mut:     &sync.Mutex{},
		screen:  screen,
		at:      At{Row: 0, Col: 0},
		atStack: []At{},
		style:   tcell.StyleDefault,
		Delay:   runeDelay,
		input:   newInput(),
	}, nil
}
