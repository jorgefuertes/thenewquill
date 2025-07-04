package adventure

import (
	"fmt"
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
)

func (a *Adventure) Export(w io.Writer) error {
	if err := a.Validate(); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "The New Quill Adventure\n"); err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(w)
	t.AppendRow(table.Row{"Title", a.Config.GetField("title")})
	t.AppendRow(table.Row{"Author", a.Config.GetField("author")})
	t.AppendRow(table.Row{"Version", fmt.Sprintf("%d", VERSION)})
	t.AppendRow(table.Row{"Date", a.Config.GetField("date")})
	t.Render()

	if err := a.DB.Export(w); err != nil {
		return err
	}

	return nil
}

func (a *Adventure) Import(r io.Reader) error {
	a.DB.Reset()

	return nil
}
