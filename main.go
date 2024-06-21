package main

import (
	"fmt"
	"os"

	repl "interpreter/repl"
)

func main() {
	fmt.Println(os.Args)
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go [filename] [--repl]")
	}
	args := os.Args[1]

	if args == "--repl" {
		repl.Start(os.Stdin, os.Stdout)
	}
}
