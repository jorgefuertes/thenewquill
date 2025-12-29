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

	n, err := a.Export(outputFilename)
	if err != nil {
		return fmt.Errorf("database export error: %w", err)
	}

	outFileInfo, err := os.Stat(outputFilename)
	if err != nil {
		return fmt.Errorf("error getting output file info: %s", err)
	}

	elapsed := time.Since(start)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Section", "Entries"})
	t.AppendRows([]table.Row{
		{"vars", fmt.Sprintf("%d", a.DB.CountRecordsByKind(kind.Variable))},
		{"vocabulary", fmt.Sprintf("%d", a.DB.CountRecordsByKind(kind.Word))},
		{"messages", fmt.Sprintf("%d", a.DB.CountRecordsByKind(kind.Message))},
		{"locations", fmt.Sprintf("%d", a.DB.CountRecordsByKind(kind.Location))},
		{"items", fmt.Sprintf("%d", a.DB.CountRecordsByKind(kind.Item))},
		{"characters", fmt.Sprintf("%d", a.DB.CountRecordsByKind(kind.Character))},
	})
	t.AppendFooter(
		table.Row{"Total", fmt.Sprintf("%d entries with %d labels", a.DB.CountRecords(), a.DB.CountLabels())},
	)
	t.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	fmt.Println()
	fmt.Printf(
		"> %s v%s\n> %s\n",
		a.Config.GetParam(config.TitleParamLabel),
		a.Config.GetParam(config.VersionParamLabel),
		a.Config.GetParam(config.AuthorParamLabel),
	)
	fmt.Printf("> Compiled in %dms\n", elapsed.Milliseconds())
	fmt.Println("> Compiler: v" + compiler.VERSION)

	fmt.Printf("> %d bytes writen to %q\n", n, outputFilename)
	fmt.Printf("> %d bytes packed size\n", outFileInfo.Size())
	fmt.Println()
	t.Render()
	fmt.Println()

	return nil
}
