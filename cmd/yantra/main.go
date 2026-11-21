package main

import (
	"fmt"
	"os"
)

const usage = `yantra - a factory that makes software factories

usage: yantra <command>

commands:
  version   print the version
`

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("yantra 0.0.0")
		return
	}
	fmt.Fprint(os.Stderr, usage)
	os.Exit(2)
}
