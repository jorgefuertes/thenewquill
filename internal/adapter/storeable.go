package adapter

import "github.com/jorgefuertes/thenewquill/internal/adventure/id"

type Storeable interface {
	GetID() id.ID
	SetID(i id.ID) Storeable
	Validate(allowNoID bool) error
}
