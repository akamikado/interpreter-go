# Interpreter in Go

# Information

Interpreter runs through a REPL

Currently supports :
- Integers
- Booleans
- Strings
    - Builtin function (len)
- Variable (dynamically typed)
- Conditionals (if else)
- First order functions
- Arrays (supports any type)
   - Builtin functions (len, first, last, tail)
- HashMaps

# How to run
- Clone the repo
- Build the project
    ```bash
        go build ./main.go`
- Run the binary (with repl flag for REPL)
    ```bash
        ./main --repl`

## TODO
- Builtin functions
    - Scanning from standard input
- Support for interpreting files
- Data structures
    - Singly Linked lists
    - Tuples
    - Stack
