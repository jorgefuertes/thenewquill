package console

import "fmt"

type ConsoleRenderer struct {
	inputOn bool
}

func New() *ConsoleRenderer {
	return new(ConsoleRenderer)
}

func (r *ConsoleRenderer) Print(s string) {
	print(s)
}

func (r *ConsoleRenderer) InputOn() {
	r.inputOn = true
}

func (r *ConsoleRenderer) InputOff() {
	r.inputOn = false
}

func (r *ConsoleRenderer) Input() string {
	if !r.inputOn {
		return ""
	}

	var input string
	fmt.Scanln(&input)

	return input
}

func (r *ConsoleRenderer) Cls() {
	fmt.Print("\033[H\033[2J")
}
