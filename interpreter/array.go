package interpreter

import (
	"glox/scanner"
)

type loxArray struct {
	size     uint
	elements []any
	token    scanner.Token
}

func newLoxArray(elements []any) *loxArray {
	return &loxArray{elements: elements, size: uint(len(elements))}
}

func (a *loxArray) validate(index uint, brackets scanner.Token) error {
	if index < 0 {
		return &Error{Token: brackets, Message: "Array indices should be positive."}
	}

	if index >= a.size {
		return &Error{Token: brackets, Message: "Array index out of bounds."}
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
	a.elements = append(a.elements, value)
	a.size += 1
	return a.elements
}
