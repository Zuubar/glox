package interpreter

import (
	"fmt"
	"glox/parser"
	"glox/scanner"
)

type loxAbstractInstance interface {
	get(name scanner.Token) (any, error)
	set(name scanner.Token, value any)
}

type loxMetaClass struct {
	stmt          parser.ClassStmt
	methods       map[string]*loxFunction
	staticMethods map[string]*loxFunction
}

func newMetaClass(class parser.ClassStmt, methods map[string]*loxFunction, staticMethods map[string]*loxFunction) *loxMetaClass {
	return &loxMetaClass{
		stmt:          class,
		methods:       methods,
		staticMethods: staticMethods,
	}
}

func (m *loxMetaClass) NewClass(superclass *loxClass) *loxClass {
	return &loxClass{
		metaClass:    m,
		superclass:   superclass,
		staticFields: make(map[string]any),
	}
}

type loxClass struct {
	metaClass    *loxMetaClass
	superclass   *loxClass
	staticFields map[string]any
}

func (c *loxClass) findMethod(name string) (*loxFunction, bool) {
	if fun, ok := c.metaClass.methods[name]; ok {
		return fun, true
	}

	if c.superclass != nil {
		return c.superclass.findMethod(name)
	}

	return nil, false
}

func (c *loxClass) get(name scanner.Token) (any, error) {
	if field, ok := c.staticFields[name.Lexeme]; ok {
		return field, nil
	}

	if staticMethod, ok := c.metaClass.staticMethods[name.Lexeme]; ok {
		return staticMethod, nil
	}

	if c.superclass != nil {
		return c.superclass.get(name)
	}

	return nil, &Error{Token: name, Message: fmt.Sprintf("Undefined property '%s'.", name.Lexeme)}
}

func (c *loxClass) set(name scanner.Token, value any) {
	c.staticFields[name.Lexeme] = value
}

func (c *loxClass) arity() int32 {
	if initializer, ok := c.metaClass.methods["init"]; ok {
		return initializer.arity()
	}

	return 0
}

func (c *loxClass) call(interpreter *Interpreter, arguments []any, token scanner.Token) (any, error) {
	instance := &loxInstance{class: c, fields: make(map[string]any)}

	if initializer, ok := c.metaClass.methods["init"]; ok {
		return initializer.bind(instance).call(interpreter, arguments, token)
	}

	return instance, nil
}

func (c *loxClass) String() string {
	return fmt.Sprintf("<class %s>", c.metaClass.stmt.Name.Lexeme)
}

type loxInstance struct {
	class  *loxClass
	fields map[string]any
}

func (i *loxInstance) get(name scanner.Token) (any, error) {
	if field, ok := i.fields[name.Lexeme]; ok {
		return field, nil
	}

	if method, ok := i.class.findMethod(name.Lexeme); ok {
		return method.bind(i), nil
	}

	return nil, &Error{Token: name, Message: fmt.Sprintf("Undefined property '%s'.", name.Lexeme)}
}

func (i *loxInstance) set(name scanner.Token, value any) {
	i.fields[name.Lexeme] = value
}

func (i *loxInstance) String() string {
	return fmt.Sprintf("<%s instance>", i.class.metaClass.stmt.Name.Lexeme)
}
