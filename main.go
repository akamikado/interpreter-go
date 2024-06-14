package main

import (
	"os"

	repl "interpreter/repl"
)

func main() {
	args := os.Args[1:]
	if args == "--repl" {
		repl.Start(os.Stdin, os.Stdout)
	}
}
