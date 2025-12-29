package adventure

import (
	"fmt"
	"time"

	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	"github.com/jorgefuertes/thenewquill/internal/adventure/item"
	"github.com/jorgefuertes/thenewquill/internal/adventure/location"
	"github.com/jorgefuertes/thenewquill/internal/adventure/message"
	"github.com/jorgefuertes/thenewquill/internal/adventure/variable"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"github.com/jorgefuertes/thenewquill/internal/database"
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
	d := database.NewDB()

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

func (a *Adventure) Export(path string) (int, error) {
	if _, err := a.Config.Set("date", fmt.Sprintf("%d", time.Now().Unix())); err != nil {
		return 0, err
	}

	_, bFile, err := a.DB.Export(path)

	return bFile, err
}

func (a *Adventure) Import(path string) error {
	return a.DB.Import(path)
}
