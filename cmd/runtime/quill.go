package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"thenewquill/internal/adventure"
	"thenewquill/internal/compiler/db"
	"thenewquill/internal/log"
	"thenewquill/internal/view"
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

	m := view.New()
	m.Println("First line")

	go func() {
		for i := 0; i < 10; i++ {
			m.Println(fmt.Sprintf("[%04d] Hello, World!", i))
			time.Sleep(time.Second * 2)
		}

		m.Println("Bye!")
	}()

	if err := m.Run(); err != nil {
		log.Fatal(err)
	}
}
