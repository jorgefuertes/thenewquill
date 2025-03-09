package renderer

type Renderer interface {
	Print(s string)
	InputOn()
	InputOff()
	Input() string
	Cls()
}
