package interpreter

import (
	"time"
)

type nativeClock struct {
}

func (n *nativeClock) arity() int32 {
	return 0
}

func (n *nativeClock) call(*Interpreter, []any) (any, error) {
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

func (n *nativeStringify) call(i *Interpreter, arguments []any) (any, error) {
	return i.stringify(arguments[0]), nil
}

func (n *nativeStringify) String() string {
	return "<native fn>"
}
