package msg

import "thenewquill/internal/compiler/section"

type MsgType int

const (
	UnknownMsg MsgType = 0
	SystemMsg  MsgType = 1
	UserMsg    MsgType = 2
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

func (t MsgType) Section() section.Section {
	switch t {
	case SystemMsg:
		return section.SysMsg
	case UserMsg:
		return section.UserMsg
	default:
		return section.None
	}
}
