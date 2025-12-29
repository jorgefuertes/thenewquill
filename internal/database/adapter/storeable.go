package adapter

type Storeable interface {
	SetID(uint32)
	GetID() uint32
	SetLabelID(uint32)
	GetLabelID() uint32
}
