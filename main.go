package main

import (
	"hek/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
