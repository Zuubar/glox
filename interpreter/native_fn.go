package interpreter

import (
	"fmt"
	"time"
)

type nativeTime struct {
}

func (n nativeTime) arity() int32 {
	return 0
}

func (n nativeTime) call(*Interpreter, []any) (any, error) {
	return float64(time.Now().UnixMilli()) / 1000, nil
}

func (n nativeTime) String() string {
	return "<native fn>"
}

type nativePrint struct {
}

func (n nativePrint) arity() int32 {
	return 1
}

func (n nativePrint) call(i *Interpreter, args []any) (any, error) {
	fmt.Print(i.stringify(args[0]))
	return nil, nil
}

func (n nativePrint) String() string {
	return "<native fn>"
}

type nativePrintLn struct {
}

func (n nativePrintLn) arity() int32 {
	return 1
}

func (n nativePrintLn) call(i *Interpreter, args []any) (any, error) {
	fmt.Println(i.stringify(args[0]))
	return nil, nil
}

func (n nativePrintLn) String() string {
	return "<native fn>"
}
