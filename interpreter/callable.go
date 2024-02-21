package interpreter

import (
	"errors"
	"fmt"
	"glox/parser"
)

type Callable interface {
	arity() int32
	call(*Interpreter, []any) (any, error)
}

type Function struct {
	funStmt parser.FunctionStmt
}

func newFunction(funDecl parser.FunctionStmt) Function {
	return Function{funStmt: funDecl}
}

func (f Function) arity() int32 {
	return int32(len(f.funStmt.Parameters))
}

func (f Function) call(interpreter *Interpreter, arguments []any) (any, error) {
	newEnv := newEnvironment(interpreter.globalEnvironment)

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

func (f Function) String() string {
	return fmt.Sprintf("<fn %s>", f.funStmt.Name.Lexeme)
}
