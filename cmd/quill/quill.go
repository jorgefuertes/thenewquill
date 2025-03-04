package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"thenewquill/internal/compiler"
	"thenewquill/internal/compiler/bin"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "quill",
		HelpName: "quill",
		Usage:    "The New Quill compiler",
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
		fmt.Println(err)
		os.Exit(1)
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

	w := bytes.NewBuffer(nil)
	w.WriteString(fmt.Sprintf("%s\n(C)%s %s %s\n", a.Config.Title, a.Config.Date, a.Config.Author, a.Config.Version))
	w.WriteString(fmt.Sprintf("Compiled by The New Quill v%s, database v%s\n", compiler.VERSION, bin.DATABASE_VERSION))
	w.WriteString("#BEGIN#")
	zw := zlib.NewWriter(w)
	if err := bin.Export(a, zw); err != nil {
		return err
	}
	zw.Close()
	w.WriteString("#END#")

	hash := sha256.Sum256(w.Bytes())
	w.Write(hash[:])

	if err := os.WriteFile(outputFilename, w.Bytes(), 0o644); err != nil {
		return err
	}

	elapsed := time.Since(start)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Section", "Entries"})
	t.AppendRows([]table.Row{
		{"vars", fmt.Sprintf("%d", a.Vars.Len())},
		{"vocabulary", fmt.Sprintf("%d", a.Words.Len())},
		{"messages", fmt.Sprintf("%d", a.Messages.Len())},
		{"locations", fmt.Sprintf("%d", a.Locations.Len())},
		{"items", fmt.Sprintf("%d", a.Items.Len())},
		{"characters", fmt.Sprintf("%d", a.Chars.Len())},
	})
	t.AppendFooter(
		table.Row{
			"Total",
			fmt.Sprintf("%d entries", a.Vars.Len()+a.Words.Len()+a.Messages.Len()+a.Locations.Len()+a.Items.Len()),
		},
	)
	t.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	fmt.Println()
	fmt.Printf("> %s v%s\n> %s\n", a.Config.Title, a.Config.Version, a.Config.Author)
	fmt.Printf("> Compiled in %dms\n", elapsed.Milliseconds())
	fmt.Println("> Compiler: v" + compiler.VERSION)
	fmt.Println("> Database: v" + bin.DATABASE_VERSION)
	fmt.Printf("> %d bytes writen to \"%s\"\n", w.Len(), outputFilename)
	fmt.Println()
	t.Render()
	fmt.Println()

	return nil
}
