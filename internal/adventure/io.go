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
	a.Reset()

	if err := a.Config.Import(d); err != nil {
		return err
	}

	if err := a.Vars.Import(d); err != nil {
		return err
	}

	a.Words.Import(d)

	if err := a.Messages.Import(d); err != nil {
		return err
	}

	if err := a.Locations.Import(d, a.Words); err != nil {
		return err
	}

	a.Chars.Import(d, a.Words, a.Locations)

	if err := a.Items.Import(d, a.Words, a.Locations, a.Chars); err != nil {
		return err
	}

	return nil
}
