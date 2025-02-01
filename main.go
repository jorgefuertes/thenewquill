package main

import (
	"fmt"
	"os"

	"thenewquill/internal/compiler"
	"thenewquill/internal/log"
)

func main() {
	a, err := compiler.Compile("test/adv/test.adv")
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	log.Info("Adventure compiled successfully")

	a.Dump()
}
