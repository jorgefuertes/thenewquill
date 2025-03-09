package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	dbFilename := ""
	if len(os.Args) > 1 {
		dbFilename = os.Args[1]
	} else {
		files, err := filepath.Glob("./*.db")
		if err != nil {
			println("❗ ERROR:", err)
		}

		for _, f := range files {
			dbFilename = f
			break
		}
	}

	if dbFilename == "" {
		fmt.Println("❗ ERROR:", ErrNoDatabase)
		os.Exit(1)
	}
}
