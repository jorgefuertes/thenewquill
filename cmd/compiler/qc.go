package main

/*
	The New Quill Compiler
	Author: queru (github.com/queru) 2025
	License: GPL-3.0
*/

import (
	"fmt"
	"os"
	"time"

	"github.com/jorgefuertes/thenewquill/internal/compiler"
	"github.com/jorgefuertes/thenewquill/internal/log"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "qc",
		HelpName: "qc",
		Usage:    "The New Quill Compiler",
		Commands: []*cli.Command{
			{
				Name:    "compile",
				Aliases: []string{"c"},
				Action:  compileAction,
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
		log.Fatal("Unexpected runtime error: %s", err)
	}
}

func compileAction(c *cli.Context) error {
	inputFilename := c.String("input")
	outputFilename := c.String("output")

	start := time.Now()

	a, err := compiler.Compile(inputFilename)
	if err != nil {
		return err
	}

	f, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal("Unexpected error clossing file: %s", err)
		}
	}()

	if err := a.Export(f); err != nil {
		return err
	}

	elapsed := time.Since(start)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Section", "Entries"})
	t.AppendRows([]table.Row{
		{"vars", fmt.Sprintf("%d", a.Variables.Count())},
		{"vocabulary", fmt.Sprintf("%d", a.Words.Count())},
		{"messages", fmt.Sprintf("%d", a.Messages.Count())},
		{"locations", fmt.Sprintf("%d", a.Locations.Count())},
		{"items", fmt.Sprintf("%d", a.Items.Count())},
		{"characters", fmt.Sprintf("%d", a.Characters.Count())},
	})
	t.AppendFooter(table.Row{"Total", fmt.Sprintf("%d entries", a.DB.Count())})
	t.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	fmt.Println()
	fmt.Printf(
		"> %s v%s\n> %s\n",
		a.Config.GetField("title"),
		a.Config.GetField("version"),
		a.Config.GetField("author"),
	)
	fmt.Printf("> Compiled in %dms\n", elapsed.Milliseconds())
	fmt.Println("> Compiler: v" + compiler.VERSION)

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	fmt.Printf("> %d bytes writen to \"%s\"\n", stat.Size(), outputFilename)
	fmt.Println()
	t.Render()
	fmt.Println()

	return nil
}
