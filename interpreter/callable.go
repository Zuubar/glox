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
	funStmt            parser.FunctionStmt
	closure            *environment
	isClassInitializer bool
	isClassGetter      bool
}

func newLoxFunction(funStmt parser.FunctionStmt, closure *environment, isClassInitializer bool, isClassGetter bool) *loxFunction {
	return &loxFunction{
		funStmt:            funStmt,
		closure:            closure,
		isClassInitializer: isClassInitializer,
		isClassGetter:      isClassGetter,
	}
}

func (f *loxFunction) bind(i loxAbstractInstance) *loxFunction {
	env := newEnvironment(f.closure)
	env.define("this", i)
	return newLoxFunction(f.funStmt, env, f.isClassInitializer, f.isClassGetter)
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
			// Handle "return;" edge case
			if f.isClassInitializer {
				this, _ := f.closure.getAt("this", 0)
				return this, nil
			}

			return returnInterrupt.Value, nil
		}

		return nil, err
	}
	// Implicit return for the class initializer
	if f.isClassInitializer {
		this, _ := f.closure.getAt("this", 0)
		return this, nil
	}

	return nil, nil
}

func (f *loxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.funStmt.Name.Lexeme)
}
