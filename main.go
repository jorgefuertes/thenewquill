package main

import (
	"fmt"
	"os"

	"thenewquill/internal/compiler"
	"thenewquill/internal/log"
)

func main() {
	a, err := compiler.Compile("internal/compiler/test/adv_files/happy/test.adv")
	if err != nil {
		fmt.Println("Compilation error:", err.Error())

		os.Exit(1)
	}

	log.Info("Adventure compiled successfully")

	a.Dump()
}
