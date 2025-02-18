package main

import (
	"fmt"
	"os"
	"thenewquill/internal/compiler"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "quill",
		HelpName: "quill",
		Usage: "The New Quill compiler",
		Commands: []*cli.Command{
			{
				Name:  "compile",
				Aliases: []string{"c"},
				Action: compileAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Usage:   "The main adventure file to compile",
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "The output file to write the compiled database",
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func compileAction(c *cli.Context) error {
	inputFilename := c.String("input")
	_ = c.String("output")

	start := time.Now()
	a, err := compiler.Compile(inputFilename)
	elapsed := time.Since(start)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Section", "Entries"})
	t.AppendRows([]table.Row{
		{"vars", fmt.Sprintf("%d", a.Vars.Len())},
		{"vocabulary", fmt.Sprintf("%d", a.Vocabulary.Len())},
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
			fmt.Sprintf("%d entries", a.Vars.Len()+a.Vocabulary.Len()+a.Messages.Len()+a.Locations.Len()+a.Items.Len()),
		},
	)
	t.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	fmt.Println()
	t.Render()
	fmt.Println()

	return nil
}
