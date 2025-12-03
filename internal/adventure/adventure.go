package adventure

import (
	"errors"
	"fmt"
	"time"

	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	"github.com/jorgefuertes/thenewquill/internal/adventure/database"
	"github.com/jorgefuertes/thenewquill/internal/adventure/item"
	"github.com/jorgefuertes/thenewquill/internal/adventure/label"
	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	"github.com/jorgefuertes/thenewquill/internal/adventure/message"
	"github.com/jorgefuertes/thenewquill/internal/adventure/variable"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
)

const VERSION = 2

type Adventure struct {
	DB         *database.DB
	Config     *config.Service
	Characters *character.Service
	Items      *item.Service
	Messages   *message.Service
	Words      *word.Service
	Locations  *location.Service
	Variables  *variable.Service
}

func New() *Adventure {
	d := database.New()
	labels := label.NewService(d)

	return &Adventure{
		DB:         d,
		Config:     config.NewService(d, labels),
		Characters: character.NewService(d, labels),
		Items:      item.NewService(d, labels),
		Messages:   message.NewService(d, labels),
		Words:      word.NewService(d, labels),
		Locations:  location.NewService(d, labels),
		Variables:  variable.NewService(d, labels),
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

func (a *Adventure) Export(path string) (int, error) {
	if _, err := a.Config.Set("date", fmt.Sprintf("%d", time.Now().Unix())); err != nil {
		return 0, err
	}

	if err := a.Validate(); err != nil {
		return 0, err
	}

	// TODO: return a.DB.Export(path)
	return 0, nil
}

func (a *Adventure) Import(path string) error {
	return nil
}
