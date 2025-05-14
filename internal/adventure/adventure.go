package adventure

import (
	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
)

const VERSION = 1

type Adventure struct {
	Config *config.Config
	Db     *db.DB
}

func New() *Adventure {
	return &Adventure{
		Config: config.New(),
		Db:     db.New(),
	}
}

func (a *Adventure) Reset() {
	a.Db.Reset()
}
