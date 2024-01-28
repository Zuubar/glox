package lox

import (
	"bufio"
	"fmt"
	"glox/lox/lox-error"
	"glox/lox/parser"
	"glox/lox/scanner"
	"os"
)

func run(source string) {
	scnr := scanner.New(source)

	tokens := scnr.Run()
	prsr := parser.New(tokens)
	ast, _ := prsr.Run()

	fmt.Printf("Scanner: %v\n", tokens)
	fmt.Printf("AST: %v\n", parser.AstPrinter{}.Print(ast))
}

func runFile(filePath string) {
	source, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	run(string(source))
}

func Repl() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		run(line)
		loxError.HadError = false
	}
}

func Run(args []string) {
	if len(args) == 0 {
		Repl()
	} else if len(args) == 1 {
		runFile(args[1])
	} else {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	}
}
