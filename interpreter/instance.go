package interpreter

import (
	"fmt"
	"glox/scanner"
)

type loxInstance struct {
	class  *loxClass
	fields map[string]any
}

func newLoxInstance(class *loxClass) *loxInstance {
	return &loxInstance{class: class, fields: make(map[string]any)}
}

func (i *loxInstance) get(name scanner.Token) (any, error) {
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

func (i *loxInstance) set(name scanner.Token, value any) {
	i.fields[name.Lexeme] = value
}

func (i *loxInstance) String() string {
	return fmt.Sprintf("<%s instance>", i.class.stmt.Name.Lexeme)
}
