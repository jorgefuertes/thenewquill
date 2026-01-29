package process

type Param byte

const (
	Str Param = iota
	Num
	Bool
	WordID
	LocID
	VarID
	ItemID
	MsgID
	TableID
	ProcID
)
