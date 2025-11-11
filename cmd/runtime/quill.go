package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"github.com/jorgefuertes/thenewquill/internal/output/console"
	"github.com/jorgefuertes/thenewquill/pkg/log"
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
		log.Fatal("Runtime Error: %s", errNoDatabase)
	}

	a := adventure.New()
	if err := a.Import(dbFilename); err != nil {
		log.Fatal("Runtime Error importing database: %s", err)
	}

	// run adventure
	o, err := console.New()
	if err != nil {
		log.Fatal("Runtime Error: %s", err)
	}

	go o.Run()
	defer o.Close()

	h, err := a.Characters.GetHuman()
	if err != nil {
		log.Fatal("Runtime Error: %s", err)
	}

	loc, err := a.Locations.Get(h.LocationID)
	if err != nil {
		log.Fatal("Runtime Error: %s", err)
	}

	o.Cls()
	o.Printf("*** %s ***\n", loc.Title)
	o.Print(loc.Description)
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

		w, err := a.Words.FirstOfAny(sl[0])
		if err != nil {
			o.Print("No te comprendo.")
			continue
		}

		if w.Is(word.Noun, "salidas") {
			o.Print("Las salidas son: ")
			for i, c := range loc.Conns {
				o.Print(a.DB.GetLabelName(c.WordID))
				if i < len(loc.Conns)-1 {
					o.Print(", ")
				} else {
					o.Println(".")
				}
			}

			for _, c := range loc.Conns {
				conLoc, err := a.Locations.Get(c.LocationID)
				if err != nil {
					log.Fatal("Runtime Error: %s", err)
				}

				o.Printf("- %s->%s: %s\n", a.DB.GetLabelName(c.WordID), a.DB.GetLabelName(c.LocationID), conLoc.Title)
			}
		}

		if loc.HasConn(w.ID) {
			id := loc.GetConn(w.ID)
			conLoc, err := a.Locations.Get(id)
			if err != nil {
				log.Fatal("Runtime Error: %s", err)
			}

			o.Printf("*** %s ***\n", conLoc.Title)
			continue
		}
	}
}
