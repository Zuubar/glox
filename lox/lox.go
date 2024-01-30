package lox

import (
	"bufio"
	"fmt"
	"glox/lox/interpreter"
	"glox/lox/parser"
	"glox/lox/scanner"
	"os"
)

func run(source string) {
	scnr := scanner.New(source)
	tokens, err := scnr.Run()

	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Printf("Scanner: %v\n", tokens)

	prsr := parser.New(tokens)
	ast, err := prsr.Run()

	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Printf("AST: %v\n", parser.AstPrinter{}.Print(ast))

	inter := interpreter.New(ast)
	result, err := inter.Run()
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Printf("Interpreted: %v\n", result)
}

func runFile(filePath string) {
	// Todo exit codes for compile-time and runtime errors
	source, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	run(string(source))
}

func repl() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		run(line)
	}
}

func Run(args []string) {
	if len(args) == 0 {
		repl()
	} else if len(args) == 1 {
		runFile(args[1])
	} else {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	}
}
