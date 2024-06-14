package main

import (
	"os"

	repl "interpreter/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
