package cmd

import (
	"bufio"
	"fmt"
	"glox/interpreter"
	"glox/parser"
	"glox/scanner"
	"os"
	"strings"
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

func run(source string) int {
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
		return 65
	}

	printer := parser.AstPrinter{}
	printDebug(fmt.Sprintf("AST: %v\n", printer.Print(ast)))

	if err := inter.Interpret(ast); err != nil {
		printErrors(err)
		return 70
	}

	return 0
}

func runFile(filePath string) {
	source, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	exitCode := run(string(source))
	os.Exit(exitCode)
}

func repl() {
	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !reader.Scan() {
			panic(reader.Err())
		}
		line := reader.Text()

		if len(line) == 0 {
			continue
		}

		if !strings.Contains(line, ";") {
			line = fmt.Sprintf("print (%s);", line)
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
