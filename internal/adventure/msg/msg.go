package msg

import "fmt"

type MsgType int

const (
	SystemMsg MsgType = iota
	UserMsg
)

func (t MsgType) String() string {
	switch t {
	case SystemMsg:
		return "system"
	case UserMsg:
		return "user"
	default:
		return "unknown"
	}
}

type Msg struct {
	Type  MsgType
	Label string
	Text  string
}

func (m Msg) String() string {
	return m.Text
}

func (m Msg) Dump() string {
	return fmt.Sprintf("[%6s] %s: %s", m.Type, m.Label, m.Text)
}

type Messages []Msg

func New() Messages {
	m := make(Messages, 0)

	return m
}

func (ms *Messages) Add(m Msg) error {
	if ms.Exists(m.Type, m.Label) {
		return ErrMsgAlreadyExists
	}

	*ms = append(*ms, m)

	return nil
}

func (ms Messages) Exists(t MsgType, label string) bool {
	for _, msg := range ms {
		if msg.Type == t && msg.Label == label {
			return true
		}
	}

	return false
}

func (ms Messages) GetText(t MsgType, label string) (string, error) {
	for _, m := range ms {
		if m.Type == t && m.Label == label {
			return m.String(), nil
		}
	}
	return "", ErrMsgNotFound
}
