package main

import (
	"glox/lox"
	"os"
)

func main() {
	lox.Run(os.Args[1:])
}
