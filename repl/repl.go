package repl

import (
	"bufio"
	"fmt"
	"io"

	l "interpreter/lexer"
)

const (
	PROMPT = ">> "
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" {
			fmt.Println("Goodbye!")
			return
		}

		l.Tokenize([]byte(line))
	}
}
