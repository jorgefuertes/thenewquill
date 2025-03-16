package adventure

import (
	"fmt"

	"thenewquill/internal/compiler/db"
)

func (a *Adventure) Export(d *db.DB) {
	d.AddHeader(fmt.Sprintf("%s v%s", a.Config.Title, a.Config.Version))
	d.AddHeader(fmt.Sprintf("By %s", a.Config.Author))
	d.AddHeader(fmt.Sprintf("Adventure Database Version %03d", VERSION))

	a.Config.Export(d)
	a.Vars.Export(d)
	a.Words.Export(d)
	a.Messages.Export(d)
	a.Locations.Export(d)
	a.Items.Export(d)
	a.Chars.Export(d)
}

func (a *Adventure) Import(d *db.DB) error {
	a = New()

	a.Config.Import(d)
	a.Words.Import(d)
	a.Messages.Import(d)
	a.Locations.Import(d, a.Words)

	if err := a.Chars.Import(d, a.Words, a.Locations); err != nil {
		return err
	}

	if err := a.Items.Import(d, a.Words, a.Locations, a.Chars); err != nil {
		return err
	}

	return nil
}
