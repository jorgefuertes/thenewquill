package msg

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
