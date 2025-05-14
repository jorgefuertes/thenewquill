package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/compiler/db"
	"github.com/jorgefuertes/thenewquill/internal/log"
	"github.com/jorgefuertes/thenewquill/internal/output/console"
)

func main() {
	dbFilename := ""
	if len(os.Args) > 1 {
		dbFilename = os.Args[1]
	} else {
		files, _ := filepath.Glob("./*.db")
		for _, f := range files {
			dbFilename = f
			break
		}
	}

	if dbFilename == "" {
		log.Fatal(errNoDatabase)
	}

	// reader from file
	r, err := os.ReadFile(dbFilename)
	if err != nil {
		log.Fatal(err)
	}

	// parse db
	d := db.New()
	err = d.Load(bytes.NewReader(r))
	if err != nil {
		log.Fatal(err)
	}

	// parse adventure
	a := adventure.New()
	err = a.Import(d)
	if err != nil {
		log.Fatal(err)
	}

	// run adventure
	o, err := console.New()
	if err != nil {
		log.Fatal(err)
	}

	go o.Run()
	defer o.Close()

	h := a.Chars.GetHuman()
	o.Cls()
	o.Printf("*** [%s] %s ***\n", h.Location.Label, h.Location.Title)
	o.Print(h.Location.Description)
	o.Println()
	time.Sleep(time.Second * 1)

	for {

		line, err := o.Input(">", time.Second*120)
		if err != nil {
			if err == console.ErrTimedOut {
				o.Print("El tiempo pasa... y no vuelve.")

				continue
			}
			if err == console.ErrCancelledByUser {
				break
			}
		}

		sl := strings.Fields(line)
		if len(sl) == 0 {
			continue
		}

		w := a.Words.First(sl[0])
		if w == nil {
			o.Print("No te comprendo.")
			continue
		}

		if w.Is("salidas") {
			o.Print("Las salidas son: ")
			for i, c := range h.Location.Conns {
				o.Print(c.Word.Label)
				if i < len(h.Location.Conns)-1 {
					o.Print(", ")
				} else {
					o.Println(".")
				}
			}

			for _, c := range h.Location.Conns {
				o.Printf("- %s->%s: %s\n", c.Word.Label, c.To.Label, c.To.Title)
			}
		}

		if h.Location.HasConn(w) {
			loc := h.Location.GetConn(w)
			o.Printf("*** %s ***\n", loc.Title)
			continue
		}
	}
}
