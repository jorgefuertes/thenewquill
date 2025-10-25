package adventure

import (
	"errors"

	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	"github.com/jorgefuertes/thenewquill/internal/adventure/db"
	"github.com/jorgefuertes/thenewquill/internal/adventure/item"
	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	"github.com/jorgefuertes/thenewquill/internal/adventure/message"
	"github.com/jorgefuertes/thenewquill/internal/adventure/variable"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
)

const VERSION = 2

type Adventure struct {
	DB         *db.DB
	Config     *config.Service
	Characters *character.Service
	Items      *item.Service
	Messages   *message.Service
	Words      *word.Service
	Locations  *location.Service
	Variables  *variable.Service
}

func New() *Adventure {
	d := db.New()

	return &Adventure{
		DB:         d,
		Config:     config.NewService(d),
		Characters: character.NewService(d),
		Items:      item.NewService(d),
		Messages:   message.NewService(d),
		Words:      word.NewService(d),
		Locations:  location.NewService(d),
		Variables:  variable.NewService(d),
	}
}

func (a *Adventure) Reset() {
	a.DB.Reset()
}

func (a *Adventure) Validate() error {
	validators := []func() error{
		a.Config.ValidateAll,
		a.Words.ValidateAll,
		a.Messages.ValidateAll,
		a.Variables.ValidateAll,
		a.Items.ValidateAll,
		a.Characters.ValidateAll,
		a.Locations.ValidateAll,
	}

	var err error

	for _, v := range validators {
		if er := v(); er != nil {
			err = errors.Join(err, er)
		}
	}

	return err
}
