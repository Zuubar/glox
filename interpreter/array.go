package interpreter

import (
	"glox/scanner"
)

type loxArray struct {
	elements []any
	token    scanner.Token
}

func newLoxArray(elements []any) *loxArray {
	return &loxArray{elements: elements}
}

func (a *loxArray) validate(index uint, brackets scanner.Token) error {
	if index < 0 {
		return &Error{Token: brackets, Message: "Array indices should be positive."}
	}

	if index >= uint(len(a.elements)) {
		return &Error{Token: brackets, Message: "Array index is out of bounds."}
	}

	return nil
}

func (a *loxArray) get(index uint) any {
	return a.elements[index]
}

func (a *loxArray) set(index uint, value any) any {
	a.elements[index] = value
	return value
}

func (a *loxArray) append(value any) any {
	return &loxArray{elements: append(a.elements, value)}
}
