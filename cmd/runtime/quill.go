package main

import (
	"fmt"
	"os"
	"time"

	"thenewquill/internal/compiler"

	"github.com/jedib0t/go-pretty/v6/table"
)

func main() {
	start := time.Now()
	a, err := compiler.Compile("internal/compiler/test/adv_files/happy/test.adv")
	elapsed := time.Since(start)
	if err != nil {
		fmt.Println("Compilation error:", err.Error())

		os.Exit(1)
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Section", "Entries"})
	t.AppendRows([]table.Row{
		{"vars", fmt.Sprintf("%d", a.Vars.Len())},
		{"vocabulary", fmt.Sprintf("%d", a.Words.Len())},
		{"messages", fmt.Sprintf("%d", a.Messages.Len())},
		{"locations", fmt.Sprintf("%d", a.Locations.Len())},
		{"items", fmt.Sprintf("%d", a.Items.Len())},
	})
	t.AppendSeparator()
	t.AppendRows([]table.Row{
		{"Compiled in", fmt.Sprintf("%dms", elapsed.Milliseconds())},
	})
	t.AppendFooter(
		table.Row{
			"Total",
			fmt.Sprintf("%d entries", a.Vars.Len()+a.Words.Len()+a.Messages.Len()+a.Locations.Len()+a.Items.Len()),
		},
	)
	t.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	t.Render()
}
