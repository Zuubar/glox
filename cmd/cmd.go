package cmd

import (
	"bufio"
	"fmt"
	"glox/interpreter"
	"glox/parser"
	"glox/scanner"
	"os"
)

var inter *interpreter.Interpreter

func printErrors(errs ...error) {
	for _, err := range errs {
		fmt.Print("\033[31m" + err.Error() + "\033[0m")
	}
}

func printDebug(message string) {
	fmt.Print("\033[33m" + message + "\033[0m")
}

func run(source string) {
	scnr := scanner.New(source)
	tokens, err := scnr.Run()

	if err != nil {
		printErrors(err)
	}
	printDebug(fmt.Sprintf("Scanner: %v\n", tokens))

	prsr := parser.New(tokens)
	ast, errs := prsr.Parse()

	if len(errs) != 0 {
		printErrors(errs...)
		return
	}

	printer := parser.AstPrinter{}
	printDebug(fmt.Sprintf("AST: %v\n", printer.Print(ast)))

	if err := inter.Interpret(ast); err != nil {
		printErrors(err)
		return
	}
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
	inter = interpreter.New()
	if len(args) == 0 {
		repl()
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	}
}
