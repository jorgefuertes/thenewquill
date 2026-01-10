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

	"github.com/jorgefuertes/thenewquill/internal/adventure/config"
	"github.com/jorgefuertes/thenewquill/internal/adventure/kind"
	"github.com/jorgefuertes/thenewquill/internal/compiler"
	"github.com/jorgefuertes/thenewquill/pkg/log"

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

	bSent, bFile, err := a.Export(outputFilename)
	if err != nil {
		return fmt.Errorf("database export error: %w", err)
	}

	elapsed := time.Since(start)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Kind", "Records"})

	for _, k := range kind.Kinds() {
		if k == kind.None || k == kind.Test || k == kind.Label {
			continue
		}

		t.AppendRow(table.Row{k.HumanName(), fmt.Sprintf("%d", a.DB.CountRecordsByKind(k))})
	}

	t.AppendFooter(
		table.Row{"Total", fmt.Sprintf("%d entries with %d labels", a.DB.CountRecords(), a.DB.CountLabels())},
	)
	t.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	fmt.Println()
	fmt.Printf(
		"> %s v%s\n> %s\n",
		a.Config.GetValueOrBlank(config.TitleParamLabel),
		a.Config.GetValueOrBlank(config.VersionParamLabel),
		a.Config.GetValueOrBlank(config.AuthorParamLabel),
	)
	fmt.Printf("> Compiled in %dms\n", elapsed.Milliseconds())
	fmt.Println("> Compiler: v" + compiler.VERSION)

	fmt.Printf("> %d bytes packed to %q\n", bFile, outputFilename)
	fmt.Printf("> %d total bytes\n", bSent)
	fmt.Println()
	t.Render()
	fmt.Println()

	return nil
}
