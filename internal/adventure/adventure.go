package adventure

import (
	"fmt"
	"io"
	"strings"

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
		a.Characters.ValidateAll,
		a.Items.ValidateAll,
		a.Messages.ValidateAll,
		a.Words.ValidateAll,
		a.Locations.ValidateAll,
		a.Variables.ValidateAll,
	}

	for _, v := range validators {
		if err := v(); err != nil {
			return err
		}
	}

	return nil
}

func (a *Adventure) Export(w io.Writer) error {
	if err := a.Validate(); err != nil {
		return err
	}

	headers := []string{
		"The New Quill Adventure",
		fmt.Sprintf("Compiler version: %d", VERSION),
		"Adventure.:" + a.Config.GetField("title"),
		"By........:" + a.Config.GetField("author"),
		"Version...:" + a.Config.GetField("version"),
		"Date......:" + a.Config.GetField("date"),
	}

	var maxLen int
	for _, h := range headers {
		if len(h) > maxLen {
			maxLen = len(h)
		}
	}

	if _, err := fmt.Fprintf(w, "+%s+\n", strings.Repeat("-", maxLen)); err != nil {
		return err
	}

	for _, h := range headers {
		if _, err := fmt.Fprintf(w, "| %s%s |\n", h, strings.Repeat(" ", maxLen-len(h))); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, "\n+%s+\n", strings.Repeat("-", maxLen)); err != nil {
		return err
	}

	if _, err := w.Write([]byte("\nBEGIN DATABASE\n")); err != nil {
		return err
	}

	if err := a.DB.Export(w); err != nil {
		return err
	}

	if _, err := w.Write([]byte("\nEND DATABASE\n")); err != nil {
		return err
	}

	return nil
}
