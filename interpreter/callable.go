package interpreter

import (
	"errors"
	"fmt"
	"glox/parser"
)

type callable interface {
	arity() int32
	call(*Interpreter, []any) (any, error)
}

type loxFunction struct {
	funStmt parser.FunctionStmt
	closure *environment
}

func newLoxFunction(funStmt parser.FunctionStmt, closure *environment) *loxFunction {
	return &loxFunction{funStmt: funStmt, closure: closure}
}

func (f *loxFunction) arity() int32 {
	return int32(len(f.funStmt.Parameters))
}

func (f *loxFunction) call(interpreter *Interpreter, arguments []any) (any, error) {
	newEnv := newEnvironment(f.closure)

	for i := 0; i < len(arguments); i++ {
		newEnv.define(f.funStmt.Parameters[i].Lexeme, arguments[i])
	}

	if _, err := interpreter.executeBlock(f.funStmt.Body, newEnv); err != nil {
		returnInterrupt := &parser.ReturnInterrupt{}
		if errors.As(err, &returnInterrupt) {
			return returnInterrupt.Value, nil
		}

		return nil, err
	}

	return nil, nil
}

func (f *loxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.funStmt.Name.Lexeme)
}

type loxClass struct {
	stmt parser.ClassStmt
}

func newClass(class parser.ClassStmt) *loxClass {
	return &loxClass{stmt: class}
}

func (c *loxClass) arity() int32 {
	return 0
}

func (c *loxClass) call(interpreter *Interpreter, arguments []any) (any, error) {
	return newLoxInstance(c), nil
}

func (c *loxClass) String() string {
	return fmt.Sprintf("<class %s>", c.stmt.Name.Lexeme)
}
