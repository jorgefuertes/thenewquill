package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jorgefuertes/thenewquill/internal/adventure"
	"github.com/jorgefuertes/thenewquill/internal/adventure/character"
	"github.com/jorgefuertes/thenewquill/internal/adventure/word"
	"github.com/jorgefuertes/thenewquill/internal/database"
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

	var h character.Character
	if err := a.DB.Query(database.NewFilter("Human", database.Equal, true)).First(&h); err != nil {
		log.Fatal("Runtime Error: human character not found")
	}

	loc, err := a.Locations.Get().WithID(h.LocationID).First()
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

		w, err := a.Words.Get().WithSynonym(sl[0]).First()
		if err != nil {
			o.Print("No te comprendo.")
			continue
		}

		if w.Is(word.Noun, "salidas") {
			o.Print("Las salidas son: ")
			for i, c := range loc.Conns {
				w, err := a.Words.Get().WithID(c.WordID).First()
				if err != nil {
					log.Fatal("Cannot get word %d: %s", c.WordID, err)
					continue
				}

				o.Print(w.Synonyms[0])
				if i < len(loc.Conns)-1 {
					o.Print(", ")
				} else {
					o.Println(".")
				}
			}

			for _, conn := range loc.Conns {
				loc, err := a.Locations.Get().WithID(conn.LocationID).First()
				if err != nil {
					log.Fatal("Cannot get location %d: %s", conn.LocationID, err)
					continue
				}

				w, err := a.Words.Get().WithID(conn.WordID).First()
				if err != nil {
					log.Fatal("Cannot get word %d: %s", conn.WordID, err)
					continue
				}

				o.Printf(
					"- %s->%s: %s\n",
					a.DB.GetLabelOrBlank(w.LabelID),
					a.DB.GetLabelOrBlank(loc.LabelID),
					loc.Title,
				)
			}
		}

		if loc.HasConn(w.ID) {
			id := loc.GetConn(w.ID)
			conLoc, err := a.Locations.Get().WithID(id).First()
			if err != nil {
				log.Fatal("Runtime Error: %s", err)
			}

			o.Printf("*** %s ***\n", conLoc.Title)
			continue
		}
	}
}
