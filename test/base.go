package test

import (
	"bytes"
	"fmt"
	"glox/interpreter"
	"glox/parser"
	"glox/resolver"
	"glox/scanner"
	"os"
	"testing"
)

type testCase struct {
	source   string
	expected string
}

func newError(t *testing.T, testCaseNum int, expected string, got string) {
	t.Fatalf("Error at the test case №%d:\nExpected: %s,\nGot %s.", testCaseNum+1, fmt.Sprintf("%q", expected), fmt.Sprintf("%q", got))
}

func testExpressions(t *testing.T, testCases []testCase) {
	for idx, tt := range testCases {
		result, err := interpret(fmt.Sprintf("print (%s);", tt.source))

		if err != nil {
			t.Fatal(err)
		}

		if len(result) > 0 {
			result = result[:len(result)-1]
		}

		if result != tt.expected {
			newError(t, idx, tt.expected, result)
		}
	}
}

func testPrograms(t *testing.T, testCases []testCase) {
	for idx, tt := range testCases {
		result, err := interpret(tt.source)

		if err != nil {
			t.Fatal(err)
		}

		if result != tt.expected {
			newError(t, idx, tt.expected, result)
		}
	}
}

func testFailingPrograms(t *testing.T, testCases []testCase) {
	for idx, tt := range testCases {
		_, err := interpret(tt.source)

		if err == nil {
			t.Fatalf("Error at the test case №%d. Error did not occur.", idx+1)
		}

		if err.Error() != tt.expected {
			newError(t, idx, tt.expected, err.Error())
		}
	}
}

func interpret(source string) (string, error) {
	_runner := func(source string, _interpreter *interpreter.Interpreter) error {
		_scanner := scanner.New(source)
		tokens, err := _scanner.Run()

		if err != nil {
			return err
		}

		_parser := parser.New(tokens)
		statements, errs := _parser.Parse()

		if len(errs) != 0 {
			return errs[0]
		}

		_resolver := resolver.New(_interpreter)
		if _, err := _resolver.Resolve(statements); err != nil {
			return err
		}

		return _interpreter.Interpret(statements)
	}

	originalStdout := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	if err := _runner(source, interpreter.New()); err != nil {
		return "", err
	}

	err := w.Close()
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		panic(err)
	}

	os.Stdout = originalStdout

	result := buf.String()

	if len(result) > 0 {
		return result, nil
	}

	return "", nil
}
