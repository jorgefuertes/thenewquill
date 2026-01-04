package compiler_error

import (
	"fmt"
	"io"
	"os"

	"github.com/jorgefuertes/thenewquill/internal/compiler/line"
	"github.com/jorgefuertes/thenewquill/internal/util"
)

const maxLineLen = 80

type output struct {
	errOutput io.Writer
	title     string
	lines     []line.Line
}

func NewOutput(title string) *output {
	return &output{title: title, lines: []line.Line{}, errOutput: os.Stderr}
}

func (o *output) SetErrOutput(w io.Writer) {
	o.errOutput = w
}

func (o *output) addNL() {
	o.lines = append(o.lines, line.Line{Num: -1, Text: ""})
}

func (o *output) addLine(n int, text string) {
	if n >= 0 {
		o.lines = append(o.lines, line.Line{Num: n, Text: util.LimitStr(text, maxLineLen-4)})

		return
	}

	if len(text) > maxLineLen-4 {
		for _, l := range util.SplitIntoLines(text, maxLineLen-4) {
			o.lines = append(o.lines, line.Line{Num: -1, Text: l})
		}

		return
	}

	o.lines = append(o.lines, line.Line{Num: n, Text: text})
}

func (o *output) Print() {
	content := errorTitleStyle.Render("*** "+o.title+" ***") + "\n"
	for _, l := range o.lines {
		c := ""
		if l.Num >= 0 {
			c += fmt.Sprintf("%05d ", l.Num)
		}

		c = decorateLine(c+l.Text) + "\n"

		content += c
	}

	fmt.Fprintln(o.errOutput, borderStyle.Render(content))
}
