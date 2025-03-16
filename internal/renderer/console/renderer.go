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

	for {
		var input string

		n, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Printf("error reading input: %v\n", err)
		}

		if n > 0 {
			return input
		}
	}
}

func (r *ConsoleRenderer) Cls() {
	fmt.Print("\033[H\033[2J")
}
