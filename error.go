package main

import (
	"flag"
	"fmt"
	"os"
)

func errorExit(msg string) {
	if msg != "" {
		fmt.Println(msg)
		fmt.Println()
	}
	flag.PrintDefaults()
	os.Exit(1)
}
