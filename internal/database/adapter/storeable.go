package adapter

import "github.com/jorgefuertes/thenewquill/internal/adventure/kind"

type Storeable interface {
	SetID(uint32)
	GetID() uint32
	SetLabelID(uint32)
	GetLabelID() uint32
	GetKind() kind.Kind
}
