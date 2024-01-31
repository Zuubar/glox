package main

import (
	"glox/cmd"
	"os"
)

func main() {
	cmd.Run(os.Args[1:])
}
