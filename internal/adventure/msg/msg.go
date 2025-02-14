package msg

import (
	"fmt"
	"strings"
)

type Msg struct {
	Type  MsgType
	Label string
	Text  string
}

func (m Msg) String() string {
	return m.Text
}

func (m Msg) Stringf(args ...any) string {
	format := strings.Replace(m.Text, "_", "%v", -1)

	return fmt.Sprintf(format, args...)
}
