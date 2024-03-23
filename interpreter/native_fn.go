package interpreter

import (
	"glox/scanner"
	"time"
)

type nativeClock struct {
}

func (n *nativeClock) arity() int32 {
	return 0
}

func (n *nativeClock) call(*Interpreter, []any, scanner.Token) (any, error) {
	return float64(time.Now().UnixMilli()) / 1000, nil
}

func (n *nativeClock) String() string {
	return "<native fn>"
}

type nativeStringify struct {
}

func (n *nativeStringify) arity() int32 {
	return 1
}

func (n *nativeStringify) call(i *Interpreter, arguments []any, _ scanner.Token) (any, error) {
	return i.stringify(arguments[0]), nil
}

func (n *nativeStringify) String() string {
	return "<native fn>"
}

type nativeAppend struct {
}

func (n *nativeAppend) arity() int32 {
	return 2
}

func (n *nativeAppend) call(i *Interpreter, arguments []any, token scanner.Token) (any, error) {
	array, ok := arguments[0].(*loxArray)
	if !ok {
		return nil, &Error{Token: token, Message: "First argument to 'append' should be an array."}
	}

	return array.append(arguments[1]), nil
}

func (n *nativeAppend) String() string {
	return "<native fn>"
}

type nativeLen struct {
}

func (n *nativeLen) arity() int32 {
	return 1
}

func (n *nativeLen) call(i *Interpreter, arguments []any, token scanner.Token) (any, error) {
	array, ok := arguments[0].(*loxArray)
	if !ok {
		return nil, &Error{Token: token, Message: "First argument to 'len' should be an array."}
	}

	return float64(len(array.elements)), nil
}

func (n *nativeLen) String() string {
	return "<native fn>"
}
