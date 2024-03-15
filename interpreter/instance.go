package interpreter

import (
	"fmt"
	"glox/parser"
	"glox/scanner"
)

type loxInstance interface {
	get(name scanner.Token) (any, error)
	set(name scanner.Token, value any)
}

type loxClass struct {
	stmt          parser.ClassStmt
	methods       map[string]*loxFunction
	staticMethods map[string]*loxFunction
	staticFields  map[string]any
}

func newClass(class parser.ClassStmt, methods map[string]*loxFunction, staticMethods map[string]*loxFunction) *loxClass {
	return &loxClass{
		stmt:          class,
		methods:       methods,
		staticMethods: staticMethods,
		staticFields:  make(map[string]any),
	}
}

func (c *loxClass) get(name scanner.Token) (any, error) {
	field, ok := c.staticFields[name.Lexeme]
	if ok {
		return field, nil
	}

	staticMethod, ok := c.staticMethods[name.Lexeme]
	if ok {
		return staticMethod, nil
	}

	return nil, &Error{Token: name, Message: fmt.Sprintf("Undefined property '%s'.", name.Lexeme)}
}

func (c *loxClass) set(name scanner.Token, value any) {
	c.staticFields[name.Lexeme] = value
}

func (c *loxClass) findMethod(name string) (*loxFunction, bool) {
	fun, ok := c.methods[name]
	return fun, ok
}

func (c *loxClass) arity() int32 {
	initializer, ok := c.methods["init"]
	if ok {
		return initializer.arity()
	}

	return 0
}

func (c *loxClass) call(interpreter *Interpreter, arguments []any) (any, error) {
	initializer, ok := c.methods["init"]
	instance := newLoxInstance(c)

	if ok {
		return initializer.bind(instance).call(interpreter, arguments)
	}

	return instance, nil
}

func (c *loxClass) String() string {
	return fmt.Sprintf("<class %s>", c.stmt.Name.Lexeme)
}

type loxClassInstance struct {
	class  *loxClass
	fields map[string]any
}

func newLoxInstance(class *loxClass) *loxClassInstance {
	return &loxClassInstance{class: class, fields: make(map[string]any)}
}

func (i *loxClassInstance) get(name scanner.Token) (any, error) {
	field, ok := i.fields[name.Lexeme]
	if ok {
		return field, nil
	}

	method, ok := i.class.findMethod(name.Lexeme)
	if ok {
		return method.bind(i), nil
	}

	return nil, &Error{Token: name, Message: fmt.Sprintf("Undefined property '%s'.", name.Lexeme)}
}

func (i *loxClassInstance) set(name scanner.Token, value any) {
	i.fields[name.Lexeme] = value
}

func (i *loxClassInstance) String() string {
	return fmt.Sprintf("<%s instance>", i.class.stmt.Name.Lexeme)
}
