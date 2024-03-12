package cmd

import (
	"bufio"
	"fmt"
	"glox/interpreter"
	"glox/parser"
	"glox/resolver"
	"glox/scanner"
	"os"
	"strings"
)

var _interpreter *interpreter.Interpreter

func printErrors(errs ...error) {
	for _, err := range errs {
		fmt.Print("\033[31m" + err.Error() + "\033[0m")
	}
}

func printWarnings(warnings ...error) {
	for _, err := range warnings {
		fmt.Print("\033[33m" + err.Error() + "\033[0m")
	}
}

func run(source string) int {
	_scanner := scanner.New(source)
	tokens, err := _scanner.Run()

	if err != nil {
		printErrors(err)
		return 63
	}

	_parser := parser.New(tokens)
	statements, errs := _parser.Parse()

	if len(errs) != 0 {
		printErrors(errs...)
		return 65
	}

	_resolver := resolver.New(_interpreter)
	if _, err := _resolver.Resolve(statements); err != nil {
		printErrors(err)
		return 67
	}

	printWarnings(_resolver.Warnings()...)

	if err := _interpreter.Interpret(statements); err != nil {
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
	// Todo: Don't use "print(expr);"
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
	_interpreter = interpreter.New()
	if len(args) == 0 {
		repl()
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		fmt.Println("Usage: glox [script]")
		os.Exit(64)
	}
}
