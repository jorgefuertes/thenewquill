package bin

import (
	"encoding/gob"
	"fmt"
	"io"
	"strings"

	"thenewquill/internal/adventure"
	"thenewquill/internal/adventure/character"
	"thenewquill/internal/adventure/msg"
	"thenewquill/internal/compiler/section"
)

func Import(r io.Reader) (*adventure.Adventure, error) {
	dto := &adventureDTO{}
	a := adventure.New()

	if err := gob.NewDecoder(r).Decode(dto); err != nil {
		return a, err
	}

	// config
	a.Config = dto.Config

	// vars
	a.Vars.SetAll(dto.Vars)

	// messages
	for _, m := range dto.Msgs {
		l, _ := dto.Labels.Get(m.ID)

		t := msg.UnknownMsg
		if l.Sec == section.UserMsg {
			t = msg.UserMsg
		} else if l.Sec == section.SysMsg {
			t = msg.SystemMsg
		}

		newMsg := msg.New(t, l.Name, m.Text)
		newMsg.SetPluralTexts(m.Plurals)

		if err := a.Messages.Set(newMsg); err != nil {
			return a, err
		}
	}

	// words
	for _, w := range dto.Words {
		l, _ := dto.Labels.Get(w.ID)

		parts := strings.Split(l.Name, "#")
		if len(parts) != 2 {
			return a, fmt.Errorf("word label #%d not found", w.ID)
		}

		_ = a.Words.Set(parts[0], w.Type, w.Synonyms...)
	}

	// locations first pass: create with no connections
	for _, l := range dto.Locations {
		label, _ := dto.Labels.Get(l.ID)
		_ = a.Locations.Set(label.Name, l.Title, l.Description)
	}

	// second pass: create the connections
	for _, l := range dto.Locations {
		from := dto.GetLocationByLabelID(l.ID, a)

		for wID, destID := range l.Conns {
			w := dto.GetWordByLabelID(wID, a)
			to := dto.GetLocationByLabelID(destID, a)
			from.SetConn(w, to)
		}
	}

	// chars
	for _, c := range dto.Chars {
		label, _ := dto.Labels.Get(c.ID)
		inLoc := dto.GetLocationByLabelID(c.LocationID, a)
		name := dto.GetWordByLabelID(c.NameID, a)
		adj := dto.GetWordByLabelID(c.AdjectiveID, a)

		newChar := character.New(label.Name, name, adj)
		newChar.Description = c.Description
		newChar.Location = inLoc
		newChar.Created = c.Created
		newChar.Human = c.Human

		if err := a.Chars.Set(newChar); err != nil {
			return a, err
		}
	}

	return a, nil
}
